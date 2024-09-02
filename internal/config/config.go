/*
Package config handles configuration for the application, including environment variables and command-line flags.

It provides functions to load environment variables from a .env file and parse command-line flags for configuration.
*/

package config

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/peterbourgon/ff/v4"
	"os"
)

var (
	Host          string // PostgreSQL host.
	DB            string // PostgreSQL database name.
	TokenPassword string // Token password for JWT.
	Port          int    // Port for the API server.
	Env           string // Environment (local, dev, prod).
	Migrations    string // Path to migration files.
	Dsn           string // PostgreSQL DSN.
)

// ParseEnv loads configuration values from a .env file into global variables.
// It uses the godotenv package to read environment variables from the specified .env file.
func ParseEnv() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file")
	}

	Host = os.Getenv("POSTGRES_HOST")
	DB = os.Getenv("POSTGRES_DB")
	TokenPassword = os.Getenv("TOKEN_PASSWORD")

	return nil
}

// ParseFlags parses command-line flags and sets the corresponding global configuration variables.
// It uses the ff package to handle flag parsing and environment variable overrides.
//
// Supported flags include:
//   - -p, --port: The port number for the API server (default: 8001).
//   - -e, --env: The environment setting (local, dev, prod) (default: local).
//   - -d, --dsn: The PostgreSQL DSN for database connection.
//   - -m, --migrations: Path to the folder containing database migration files.
//
// If an invalid environment value is provided, an error is returned.
func ParseFlags(args []string) error {
	flagSet := ff.NewFlagSet("API")

	flagSet.IntVar(&Port, 'p', "port", 8001, "API server port")
	flagSet.StringVar(&Env, 'e', "env", "local", "Environment (local|dev|prod)")
	flagSet.StringVar(&Dsn, 'd', "dsn", "", "PostgreSQL DSN")
	flagSet.StringVar(&Migrations, 'm', "migrations", "", "Path to migration files folder. If not provided, migrations do not apply")

	err := ff.Parse(flagSet, args, ff.WithEnvVarPrefix("APP"))
	if err != nil {
		return err
	}

	if Env != "local" && Env != "dev" && Env != "prod" {
		return errors.New("invalid env value")
	}

	return nil
}
