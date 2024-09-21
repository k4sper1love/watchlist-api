package rest

import (
	"context"
	"errors"
	"fmt"
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
	httpAddr := fmt.Sprintf(":%d", config.Port)

	// Create a new HTTP server with configured address and timeouts.
	httpServer := &http.Server{
		Addr:         httpAddr,
		Handler:      route(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Create a new HTTPS server with configured address and timeouts.
	httpsServer := &http.Server{
		Addr:         ":443",
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
		signal.Notify(quit, os.Interrupt, os.Kill)

		// Wait for a termination signal.
		s := <-quit

		sl.Log.Debug("caught signal", slog.String("signal", s.String()))

		// Create a context with timeout for the shutdown process.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to shut down the servers gracefully.
		shutdownErr <- httpServer.Shutdown(ctx)
		shutdownErr <- httpsServer.Shutdown(ctx)
	}()

	// Getting the value of the USE_HTTPS environment variable
	useHTTPS := os.Getenv("USE_HTTPS")

	// Checking if HTTPS is enabled
	if useHTTPS == "true" {
		// Setting up autocert for automatic HTTPS
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(os.Getenv("SERVER_HOST")),
			Cache:      autocert.DirCache("certs"),
		}

		// Launching an HTTPS server in goroutine
		go func() {
			sl.Log.Info("starting HTTPS server", slog.String("address", httpsServer.Addr))

			err := httpsServer.Serve(m.Listener())
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				sl.Log.Error("https server error", slog.Any("error", err))
			}
		}()

		// HTTP handler for redirecting to HTTPS
		httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + os.Getenv("SERVER_HOST") + ":443" + r.RequestURI
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})

	}

	sl.Log.Info("starting HTTP server", slog.String("address", httpServer.Addr))
	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		sl.Log.Error("http server error", slog.Any("error", err))
		return err
	}

	// Wait for server shutdown to complete and handle any shutdown errors.
	err = <-shutdownErr
	if err != nil {
		sl.Log.Error("shutdown error", slog.Any("error", err))
		return err
	}

	sl.Log.Info("stopped server gracefully")

	return nil
}
