package watchlist

import (
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/transport/rest"
	"log"
)

func Run() {
	log.Println("Run initialization whole app")

	log.Print("Loading .env file")
	err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Print("Connection to database")
	db := postgres.ConnectPostgres()
	if db == nil {
		log.Fatal("Cannot connect to PostgreSQL")
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Print("Loading server")
	server := rest.LoadServer()
	if server == nil {
		log.Fatal("Error loading server")
	}

	log.Println("Run server on", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}
