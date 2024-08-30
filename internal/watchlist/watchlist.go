package watchlist

import (
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/transport/rest"
	"log"
)

func Run() {
	log.Println("run initialization whole app")

	log.Print("loading .env file")
	err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	log.Print("connection to database")
	db := postgres.ConnectPostgres()
	if db == nil {
		log.Fatal("cannot connect to PostgreSQL")
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
