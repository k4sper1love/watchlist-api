// Package watchlist initializes and starts the application.
//
// It handles the following tasks:
// 1. Sets up logging with configurable formats based on the environment.
// 2. Loads configuration from environment variables and command-line flags.
// 3. Establishes a connection to the PostgreSQL database.
// 4. Starts the REST API server.
//
// The Run function is the entry point for starting the application and manages the overall setup and execution flow.
package watchlist

import (
	"github.com/k4sper1love/watchlist-api/api"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/transport/rest"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/metrics"
	"github.com/k4sper1love/watchlist-api/pkg/version"
	"log/slog"
)

// Run initializes and starts the application, handling configuration,
// logging, database connection, and server startup.
func Run(args []string) error {
	slog.Info("starting application")

	if err := config.ParseFlags(args[1:]); err != nil {
		return err
	}

	// Reconfigure logger based on the environment.
	sl.Init(config.Env)
	api.SwaggerInfo.Version = version.GetVersion()

	// Open database connection
	postgres.OpenDB()

	defer postgres.CloseDB()

	// Start the REST server.
	metrics.InitUptime()
	return rest.Serve()
}
