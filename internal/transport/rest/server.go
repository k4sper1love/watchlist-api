package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Serve starts the HTTP server and handles graceful shutdown on receiving termination signals.
// It listens for incoming HTTP requests and routes them using the configured route handler.
// Returns an error if the server encounters a problem or fails to shut down gracefully.
func Serve() error {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// Create a new HTTP server with configured address and timeouts.
	server := &http.Server{
		Addr:         address,
		Handler:      route(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Channel to receive errors during server shutdown.
	shutdownErr := make(chan error)

	go func() {
		// Channel to receive termination signals.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Wait for a termination signal.
		s := <-quit

		sl.Log.Debug("caught signal", slog.String("signal", s.String()))

		// Create a context with timeout for the shutdown process.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to shut down the server gracefully.
		shutdownErr <- server.Shutdown(ctx)
	}()

	sl.Log.Info("starting server", slog.String("address", server.Addr))

	// Start the server and listen for incoming requests.
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		sl.Log.Error("server error", slog.Any("error", err))
		return err
	}

	// Wait for server shutdown to complete and handle any shutdown errors.
	err = <-shutdownErr
	if err != nil {
		sl.Log.Error("shutdown error", slog.Any("error", err))
		return err
	}

	sl.Log.Info("stopped server on", slog.String("address", server.Addr))

	return nil
}
