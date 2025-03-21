// Package config handles configuration for the application from command-line flags.
//
// It provides functions to parse command-line flags for configuration.
package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/peterbourgon/ff/v4"
	"log/slog"
	"os"
)

var (
	Env            string // Environment (local, dev, prod).
	Migrations     string // Path to migration files.
	Dsn            string // PostgreSQL Data Source Name for database connection.
	Port           int    // Port for the API server.
	JWTSecret      string // Secret password for creating JWT tokens.
	TelegramSecret string // Secret password for checking verification token
)

// ParseFlags parses command-line flags and sets the corresponding global configuration variables.
// It uses the ff package to handle flag parsing and environment variable overrides.
//
// Supported flags include:
//   - -p, --port: The port number for the API server (default: 8001).
//   - -e, --env: The environment setting (local, dev, prod) (default: local).
//   - -m, --migrations: Path to the folder containing database migration files.
//   - -s, --secret: The secret password for creating JWT tokens.
//   - -t, --telegram: The secret password for checking verification token
func ParseFlags(args []string) error {
	// Create a new flag set for the API configuration
	flagSet := ff.NewFlagSet("API Configuration")

	// Define command-line flags and their default values
	flagSet.IntVar(&Port, 'p', "port", 8001, "API server port")
	flagSet.StringVar(&Env, 'e', "env", "local", "Environment (local|dev|prod)")
	flagSet.StringVar(&Migrations, 'm', "migrations", "", "Path to migration files folder. If not provided, migrations do not apply")
	flagSet.StringVar(&JWTSecret, 's', "secret", "secretPass", "Secret password for creating JWT tokens")
	flagSet.StringVar(&TelegramSecret, 't', "telegram", "secretPassq", "Secret password for checking verification token")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		slog.Debug("no .env file found")
	}

	// Parse flags and environment variables
	if err := ff.Parse(flagSet, args, ff.WithEnvVarPrefix("APP")); err != nil {
		slog.Error("error parsing flags", slog.Any("error", err))
		return err
	}

	//Compose the PostgreSQL DSN from environment variables
	composePostgresDSN()

	// Validate the application environment
	if !isValidEnv(Env) {
		slog.Warn("invalid environment value; defaulting to 'local'", slog.Any("env", Env))
		Env = "local"
	}

	slog.Debug("successfully parsed flags")
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
