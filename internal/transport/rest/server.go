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

// Serve initializes and starts the HTTP(S) server, handling both HTTP requests
// and graceful shutdown when termination signals are received. It dynamically
// decides whether to start an HTTP or HTTPS server based on the USE_HTTPS
// environment variable. In HTTPS mode, it also sets up a server for HTTP to HTTPS redirection.
func Serve() error {
	useHTTPS := os.Getenv("USE_HTTPS")
	portHTTP := fmt.Sprint(config.Port)
	portHTTPS := "443"

	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}

	// Create servers
	serverHTTP := newServer(":" + portHTTP)
	serverHTTPS := newServer(":" + portHTTPS)

	// Setup autocert for automatic HTTPS certificate management.
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),
		Cache:      autocert.DirCache("certs"),
	}
	serverHTTPS.TLSConfig = m.TLSConfig() // Configures the HTTPS server to use autocert

	// Channel for graceful shutdown
	shutdownErr := make(chan error)
	go gracefulShutdown(serverHTTP, serverHTTPS, shutdownErr)

	if useHTTPS == "true" {
		go startHTTPS(serverHTTPS, host)

		serverHTTP.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			redirectToHTTPS(w, r, host)
		})

		api.SwaggerInfo.Host = host
		api.SwaggerInfo.Schemes = []string{"https"}
	} else {
		api.SwaggerInfo.Host = host + ":" + portHTTP
		api.SwaggerInfo.Schemes = []string{"http"}
	}

	sl.Log.Info("starting HTTP server", slog.String("address", "http://"+host+serverHTTP.Addr))

	if err := serverHTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		sl.Log.Error("HTTP server error", slog.Any("error", err))
		return err
	}

	if err := <-shutdownErr; err != nil {
		sl.Log.Warn("shutdown error", slog.Any("error", err))
		return err
	}

	sl.Log.Info("stopped server gracefully")
	return nil
}

func newServer(addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      route(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}
}

func gracefulShutdown(serverHTTP, serverHTTPS *http.Server, shutdownErr chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	s := <-quit
	sl.Log.Debug("caught signal", slog.String("signal", s.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdownErr <- serverHTTP.Shutdown(ctx)
	shutdownErr <- serverHTTPS.Shutdown(ctx)
}

func startHTTPS(serverHTTPS *http.Server, host string) {
	sl.Log.Info("starting HTTPS server", slog.String("address", "https://"+host+serverHTTPS.Addr))

	err := serverHTTPS.ListenAndServeTLS("", "")

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		sl.Log.Error("HTTPS server error", slog.Any("error", err))
	}
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request, host string) {
	target := "https://" + host + r.RequestURI
	http.Redirect(w, r, target, http.StatusMovedPermanently)

	sl.Log.Debug(
		"redirecting to HTTPS",
		slog.String("original_url", r.URL.String()),
		slog.String("target_url", "https://"+host+r.RequestURI),
		slog.String("from", r.RemoteAddr),
	)
}
