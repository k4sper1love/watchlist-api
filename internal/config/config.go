/*
Package config handles configuration for the application from command-line flags.

It provides functions to parse command-line flags for configuration.
*/

package config

import (
	"github.com/joho/godotenv"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/peterbourgon/ff/v4"
	"log/slog"
)

var (
	TokenPass  string // Token password for JWT.
	Port       int    // Port for the API server.
	Env        string // Environment (local, dev, prod).
	Migrations string // Path to migration files.
	Dsn        string // PostgreSQL DSN.
)

// ParseFlags parses command-line flags and sets the corresponding global configuration variables.
// It uses the ff package to handle flag parsing and environment variable overrides.
//
// Supported flags include:
//   - -p, --port: The port number for the API server (default: 8001).
//   - -e, --env: The environment setting (local, dev, prod) (default: local).
//   - -d, --dsn: The PostgreSQL DSN for database connection.
//   - -m, --migrations: Path to the folder containing database migration files.
//   - -s, --secret: The secret password for creating JWT tokens.
//
// If an invalid environment value is provided, an error is returned.
func ParseFlags(args []string) error {
	flagSet := ff.NewFlagSet("API")

	flagSet.IntVar(&Port, 'p', "port", 8001, "API server port")
	flagSet.StringVar(&Env, 'e', "env", "local", "Environment (local|dev|prod)")
	flagSet.StringVar(&Dsn, 'd', "dsn", "", "PostgreSQL DSN")
	flagSet.StringVar(&Migrations, 'm', "migrations", "", "Path to migration files folder. If not provided, migrations do not apply")
	flagSet.StringVar(&TokenPass, 's', "secret", "secretPass", "Secret password for creating JWT tokens")

	err := godotenv.Load()
	if err != nil {
		sl.Log.Debug("not found .env file", slog.Any("error", err))
	}

	err = ff.Parse(flagSet, args, ff.WithEnvVarPrefix("APP"))
	if err != nil {
		sl.Log.Error("Error parsing flags", slog.Any("error", err))
		return err
	}

	if Env != "local" && Env != "dev" && Env != "prod" {
		sl.Log.Warn("Invalid environment value; defaulting to 'local'", slog.Any("env", Env))
		Env = "local"
	}

	sl.Log.Debug("Parsed flags successfully", slog.String("env", Env), slog.Int("port", Port))
	return nil
}
