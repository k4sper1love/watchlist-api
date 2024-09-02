package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/vcs"
	"net/http"
)

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	message := envelope{
		"status": "enabled",
		"system_info": envelope{
			"environment": config.Env,
			"version":     vcs.GetVersion(),
		},
	}

	writeJSON(w, r, http.StatusOK, message)
}
