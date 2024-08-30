package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/vcs"
	"log"
	"net/http"
)

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("healthcheckHandler serving:", r.URL.Path, r.Host)

	message := envelope{
		"status": "enabled",
		"system_info": envelope{
			"environment": "none",
			"version":     vcs.Version(),
		},
	}

	writeJSON(w, r, http.StatusOK, message)
}
