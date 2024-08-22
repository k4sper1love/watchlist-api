package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var (
	Hostname        string
	Port            int
	Username        string
	Password        string
	Database        string
	ApplicationPort int
)

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file")
	}

	Hostname = os.Getenv("POSTGRES_HOSTNAME")
	Username = os.Getenv("POSTGRES_USERNAME")
	Password = os.Getenv("POSTGRES_PASSWORD")
	Database = os.Getenv("POSTGRES_DATABASE")

	Port, err = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return err
	}

	ApplicationPort, err = strconv.Atoi(os.Getenv("APPLICATION_PORT"))
	if err != nil {
		return err
	}

	return nil
}
