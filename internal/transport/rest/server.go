package rest

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"net/http"
	"time"
)

func LoadServer() *http.Server {
	address := fmt.Sprintf("%s:%d", config.Hostname, config.ApplicationPort)

	server := &http.Server{
		Addr:         address,
		Handler:      route(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	return server
}
