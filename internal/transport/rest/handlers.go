// Package rest provides HTTP handlers for managing and retrieving information related to the REST API.
//
// This package includes handlers for adding, retrieving, updating, and deleting users, films,
// collections, and collection-films, as well as for checking the health of the API.
//
// The handlers use a custom logger for logging and interact with the database and other internal
// components to perform various operations related to users, films, collections, and permissions.

package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/vcs"
	"net/http"
)

// HealthCheck godoc
// @Summary Check API status
// @Description Check the API status. Returns status and system information.
// @Tags monitoring
// @Accept json
// @Produce json
// @Success 200 {object} swagger.HealthCheckResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Router /healthcheck [get]
func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Construct a response message with status and system information.
	message := envelope{
		"status": "enabled",
		"system_info": envelope{
			"environment": config.Env,       // Current environment (e.g., dev, prod).
			"version":     vcs.GetVersion(), // Current version of the API.
		},
	}

	writeJSON(w, r, http.StatusOK, message)
}
