package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var (
	Host            string
	Port            int
	User            string
	Pass            string
	Database        string
	ApplicationPort int
)

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file")
	}

	Host = os.Getenv("POSTGRES_HOST")
	User = os.Getenv("POSTGRES_USER")
	Pass = os.Getenv("POSTGRES_PASS")
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
