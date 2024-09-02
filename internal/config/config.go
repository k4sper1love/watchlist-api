package config

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/peterbourgon/ff/v4"
	"os"
)

var (
	Host          string
	DB            string
	TokenPassword string
	Port          int
	Env           string
	Migrations    string
	Dsn           string
)

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

func ParseFlags(args []string) error {
	flagSet := ff.NewFlagSet("API")

	flagSet.IntVar(&Port, 'p', "port", 8001, "API server port")
	flagSet.StringVar(&Env, 'e', "env", "dev", "Environment (local|dev|prod)")
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
