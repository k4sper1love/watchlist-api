package rest

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/version"
	"net/http"
	"time"
)

var (
	message    = "Welcome to the Watchlist API"
	apiVersion = "v1"
	createdBy  = "github.com/k4sper1love"
	license    = "MIT"

	availableEndpoints = []Endpoint{
		{Path: "/swagger/index.html", Description: "API documentation"},
		{Path: "/api/v1/healthcheck", Description: "Server status"},
		{Path: "/api/v1/auth", Description: "Authorization"},
		{Path: "/api/v1/user", Description: "Profile Management"},
		{Path: "/api/v1/films", Description: "Films Management"},
		{Path: "/api/v1/collections", Description: "Collections Management"},
	}
)

type Endpoint struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

type APIResponse struct {
	Message            string     `json:"message"`
	Version            string     `json:"version"`
	APIVersion         string     `json:"api_version"`
	Timestamp          string     `json:"timestamp"`
	AvailableEndpoints []Endpoint `json:"available_endpoints"`
	CreatedBy          string     `json:"created_by"`
	License            string     `json:"license"`
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	response := APIResponse{
		Message:            message,
		Version:            version.GetVersion(),
		APIVersion:         apiVersion,
		Timestamp:          time.Now().Format(time.RFC3339),
		AvailableEndpoints: availableEndpoints,
		CreatedBy:          createdBy,
		License:            license,
	}

	writeJSON(w, r, http.StatusOK, envelope{"api": response})
}
