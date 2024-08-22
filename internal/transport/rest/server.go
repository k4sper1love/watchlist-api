package rest

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"net/http"
	"time"
)

var Address string

var Server *http.Server

func LoadServer() error {
	Address = fmt.Sprintf("%s:%d", config.Hostname, config.ApplicationPort)

	Server = &http.Server{
		Addr:         Address,
		Handler:      route(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	return nil
}
