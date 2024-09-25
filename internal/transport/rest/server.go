package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/k4sper1love/watchlist-api/api"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"golang.org/x/crypto/acme/autocert"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Serve initializes and starts the HTTP(S) server based on the USE_HTTPS environment variable.
// Handles graceful shutdown when receiving termination signals.
func Serve() error {
	useHTTPS := os.Getenv("USE_HTTPS") == "true"
	port := fmt.Sprint(config.Port)
	host := getServerHost()

	server := newServer(port)

	shutdownErr := make(chan error)
	go handleGracefulShutdown(server, shutdownErr)

	if useHTTPS {
		if err := startHTTPS(server, host, port); err != nil {
			return err
		}
	} else {
		if err := startHTTP(server, host, port); err != nil {
			return err
		}
	}

	if err := <-shutdownErr; err != nil {
		sl.Log.Warn("shutdown error", slog.Any("error", err))
		return err
	}

	sl.Log.Info("stopped server gracefully")
	return nil
}

// getServerHost returns the server host or defaults to "localhost".
func getServerHost() string {
	if host := os.Getenv("SERVER_HOST"); host != "" {
		return host
	}
	return "localhost"
}

// newServer creates a new HTTP server with common configurations.
func newServer(addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      route(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}
}

// startHTTP configures and starts the HTTP server.
func startHTTP(server *http.Server, host, port string) error {
	api.SwaggerInfo.Host = host + port
	api.SwaggerInfo.Schemes = []string{"http"}

	sl.Log.Info("starting HTTP server", slog.String("address", "http://"+host+server.Addr))

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		sl.Log.Error("HTTP server error", slog.Any("error", err))
		return err
	}
	return nil
}

// startHTTPS configures and starts the HTTPS server with automatic certificate management.
func startHTTPS(server *http.Server, host, port string) error {
	api.SwaggerInfo.Host = host
	api.SwaggerInfo.Schemes = []string{"https"}

	configureSSL(server, host)

	sl.Log.Info("starting HTTPS server", slog.String("address", "https://"+host+server.Addr+port))

	if err := server.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
		sl.Log.Error("HTTPS server error", slog.Any("error", err))
		return err
	}
	return nil
}

// configureSSL sets up SSL/TLS configuration using autocert for automatic certificate management.
func configureSSL(server *http.Server, host string) {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),
		Cache:      autocert.DirCache("certs"),
	}

	server.TLSConfig = m.TLSConfig()
}

// handleGracefulShutdown listens for termination signals and gracefully shuts down the server.
func handleGracefulShutdown(server *http.Server, shutdownErr chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	s := <-quit
	sl.Log.Debug("caught signal", slog.String("signal", s.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdownErr <- server.Shutdown(ctx)
}
