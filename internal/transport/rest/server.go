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
	serverHost := os.Getenv("SERVER_HOST")
	httpPort := fmt.Sprint(config.Port)
	httpsPort := "443"

	// Create a new HTTP server with configured address and timeouts.
	httpServer := &http.Server{
		Addr:         ":" + httpPort,
		Handler:      route(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Create a new HTTPS server with configured address and timeouts.
	httpsServer := &http.Server{
		Addr:         ":" + httpsPort,
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

	//Checking if HTTPS is enabled
	if useHTTPS == "true" {
		// Setting up autocert for automatic HTTPS
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(serverHost),
			Cache:      autocert.DirCache("/certs"),
		}

		httpsServer.TLSConfig = m.TLSConfig()

		// Launching an HTTPS server in goroutine
		go func() {
			sl.Log.Info("starting HTTPS server", slog.String("address", "https://"+serverHost+httpsServer.Addr))

			err := httpsServer.ListenAndServeTLS("", "")
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				sl.Log.Error("https server error", slog.Any("error", err))
			}
		}()

		// HTTP handler for redirecting to HTTPS
		httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + serverHost + r.RequestURI
			http.Redirect(w, r, target, http.StatusMovedPermanently)
			sl.Log.Info(
				"Redirecting to HTTPS",
				slog.String("original_url", r.URL.String()),
				slog.String("target_url", "https://"+serverHost+r.RequestURI),
				slog.String("from", r.RemoteAddr),
			)
		})

		// If HTTPS is enabled, set the host and scheme for Swagger to use HTTPS.
		api.SwaggerInfo.Host = fmt.Sprintf("%s", serverHost)
		api.SwaggerInfo.Schemes = []string{"https"}
	} else {
		// If HTTPS is not enabled, set the host and scheme for Swagger to use HTTP.
		api.SwaggerInfo.Host = fmt.Sprintf("%s:%s", serverHost, httpPort)
		api.SwaggerInfo.Schemes = []string{"http"}
	}

	sl.Log.Info("starting HTTP server", slog.String("address", "http://"+serverHost+httpServer.Addr))

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		sl.Log.Error("http server error", slog.Any("error", err))
		return err
	}

	// Wait for server shutdown to complete and handle any shutdown errors.
	err = <-shutdownErr
	if err != nil {
		sl.Log.Warn("shutdown error", slog.Any("error", err))
		return err
	}

	sl.Log.Info("stopped server gracefully")

	return nil
}
