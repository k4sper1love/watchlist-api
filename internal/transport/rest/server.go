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

func Serve() error {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	server := &http.Server{
		Addr:         address,
		Handler:      route(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	shutdownErr := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		sl.Log.Debug("caught signal", slog.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr <- server.Shutdown(ctx)
	}()

	sl.Log.Info("starting server", slog.String("address", server.Addr))

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		sl.Log.Error("server error", slog.Any("error", err))
		return err
	}

	err = <-shutdownErr
	if err != nil {
		sl.Log.Error("shutdown error", slog.Any("error", err))
		return err
	}

	sl.Log.Info("stopped server on", slog.String("address", server.Addr))

	return nil
}
