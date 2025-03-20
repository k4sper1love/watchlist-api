/*
Package main initializes and starts the Watchlist API application.

The `main` function calls `watchlist.Run` with command-line arguments to start the application.
*/

package main

import (
	"github.com/k4sper1love/watchlist-api/internal/watchlist"
	"log/slog"
	"os"
)

// @title Watchlist API
// @description This is a REST API for saving films you want to watch.
// @BasePath /api/v1

// @securityDefinitions.apiKey JWTAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: 'Authorization: Bearer {token}'

func main() {
	if err := watchlist.Run(os.Args); err != nil {
		slog.Error("application terminated due to an error")
		os.Exit(1)
	}
}
