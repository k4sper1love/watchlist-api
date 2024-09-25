/*
Package config handles configuration for the application from command-line flags.

It provides functions to parse command-line flags for configuration.
*/

package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/peterbourgon/ff/v4"
	"log/slog"
	"os"
)

var (
	JwtSecret  string   // Secret password for creating JWT tokens.
	Env        string   // Environment (local, dev, prod).
	Migrations string   // Path to migration files.
	Dsn        string   // PostgreSQL Data Source Name for database connection.
	Port       = "8001" // Port for the API server.
)

// ParseFlags parses command-line flags and sets the corresponding global configuration variables.
// It uses the ff package to handle flag parsing and environment variable overrides.
//
// Supported flags include:
//   - -e, --env: The environment setting (local, dev, prod) (default: local).
//   - -m, --migrations: Path to the folder containing database migration files.
//   - -s, --secret: The secret password for creating JWT tokens.
func ParseFlags(args []string) error {
	// Create a new flag set for the API configuration
	flagSet := ff.NewFlagSet("API Configuration")

	// Define command-line flags and their default values
	flagSet.StringVar(&Env, 'e', "env", "local", "Environment (local|dev|prod)")
	flagSet.StringVar(&Migrations, 'm', "migrations", "", "Path to migration files folder. If not provided, migrations do not apply")
	flagSet.StringVar(&JwtSecret, 's', "secret", "secretPass", "Secret password for creating JWT tokens")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		sl.Log.Debug("no .env file found")
	}

	// Parse flags and environment variables
	if err := ff.Parse(flagSet, args, ff.WithEnvVarPrefix("APP")); err != nil {
		sl.Log.Error("error parsing flags", slog.Any("error", err))
		return err
	}

	//Compose the PostgreSQL DSN from environment variables
	composePostgresDSN()

	// Validate the application environment
	if !isValidEnv(Env) {
		sl.Log.Warn("invalid environment value; defaulting to 'local'", slog.Any("env", Env))
		Env = "local"
	}

	sl.Log.Debug("parsed flags successfully")
	return nil
}

func composePostgresDSN() {
	Dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}

func isValidEnv(env string) bool {
	validEnvs := map[string]struct{}{
		"local": {},
		"dev":   {},
		"prod":  {},
	}
	_, exists := validEnvs[env]
	return exists
}
