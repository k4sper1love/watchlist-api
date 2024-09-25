package rest

import (
	"context"
	"errors"
	"github.com/k4sper1love/watchlist-api/api"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Serve initializes and starts the HTTP(S) server based on the USE_HTTPS environment variable.
// Handles graceful shutdown when receiving termination signals.
func Serve() error {
	host := getServerHost()
	server := newServer(config.Port)

	shutdownErr := make(chan error)
	go handleGracefulShutdown(server, shutdownErr)

	if err := startHTTP(server, host); err != nil {
		return err
	}

	if err := <-shutdownErr; err != nil {
		sl.Log.Warn("shutdown error", slog.Any("error", err))
		return err
	}

	sl.Log.Info("stopped server gracefully")
	return nil
}

// newServer creates a new HTTP server with common configurations.
func newServer(port string) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      route(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}
}

// startHTTP configures and starts the HTTP server.
func startHTTP(server *http.Server, host string) error {
	api.SwaggerInfo.Host = host + server.Addr
	api.SwaggerInfo.Schemes = []string{"http"}

	sl.Log.Info("starting HTTP server", slog.String("address", "http://"+host+server.Addr))

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		sl.Log.Error("HTTP server error", slog.Any("error", err))
		return err
	}
	return nil
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

// getServerHost returns the server host or defaults to "localhost".
func getServerHost() string {
	if host := os.Getenv("SERVER_HOST"); host != "" {
		return host
	}
	return "localhost"
}
