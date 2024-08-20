package watchlist

import (
	"github.com/k4sper1love/watchlist-api/internal/config"
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

	log.Print("Loading server")
	err = rest.LoadServer()
	if err != nil {
		log.Fatal("Error loading server")
	}

	log.Println("Run server on", rest.Address)
	err = rest.Server.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}
