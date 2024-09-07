/*
Package watchlist initializes and starts the application.

It handles the following tasks:
1. Sets up logging with configurable formats based on the environment.
2. Loads configuration from environment variables and command-line flags.
3. Establishes a connection to the PostgreSQL database.
4. Starts the REST API server.

The Run function is the entry point for starting the application and manages the overall setup and execution flow.
*/

package watchlist

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/docs"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/transport/rest"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/vcs"
	"log/slog"
	"os"
)

// Run initializes and starts the application, handling configuration,
// logging, database connection, and server startup.
func Run(args []string) {
	// Set up logging for the application, outputting to standard output.
	sl.SetupLogger("local", os.Stdout)
	sl.Log.Info("starting application")

	// Load environment variables from the .env file.
	err := config.ParseEnv()
	if err != nil {
		sl.Log.Error("failed to load .env file", slog.Any("error", err))
		os.Exit(1) // Exit if loading environment variables fails.
	}

	sl.Log.Debug("environment variables parsed successfully")

	// Parse command-line flags.
	err = config.ParseFlags(args[1:])
	if err != nil {
		sl.Log.Error("failed to load flags", slog.Any("error", err))
		os.Exit(1) // Exit if flag parsing fails.
	}

	sl.Log.Debug("command-line flags parsed successfully")

	// Reconfigure logger based on the environment.
	sl.SetupLogger(config.Env, os.Stdout)

	// Open a connection to the database.
	db, err := postgres.OpenDB()
	if err != nil {
		os.Exit(1) // Exit if database connection fails.
	}

	// Ensure the database connection is closed when the function exits.
	defer func() {
		if err := db.Close(); err != nil {
			sl.Log.Error("failed to close database connection", slog.Any("error", err))
			os.Exit(1) // Exit if closing the database connection fails.
		}
		sl.Log.Info("database connection closed")
	}()

	// Configure Swagger documentation.
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", config.Host, config.Port)
	docs.SwaggerInfo.Version = vcs.GetVersion()

	// Start the REST server.
	err = rest.Serve()
	if err != nil {
		os.Exit(1) // Exit if server startup fails.
	}
}
