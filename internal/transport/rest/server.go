package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"log"
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

		log.Println("caught signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr <- server.Shutdown(ctx)
	}()

	log.Println("starting server on", server.Addr)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	log.Println("stopped server on", server.Addr)

	return nil

}
