package watchlist

import (
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/transport/rest"
	"log"
)

func Run(args []string) {
	err := config.ParseEnv()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	err = config.ParseFlags(args[1:])
	if err != nil {
		log.Fatalf("error loading flags: %v", err)
	}

	db, err := postgres.OpenDB()
	if err != nil {
		log.Fatalf("cannot connect to PostgreSQL or apply migrations: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err = rest.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
