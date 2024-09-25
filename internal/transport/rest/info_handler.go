package rest

import (
	"github.com/k4sper1love/watchlist-api/pkg/version"
	"net/http"
	"time"
)

var (
	apiVersion = "v1"
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
	response := APIResponse{
		Message:            "Welcome to the Watchlist API",
		Version:            version.GetVersion(),
		APIVersion:         apiVersion,
		Timestamp:          time.Now().Format(time.RFC3339),
		AvailableEndpoints: availableEndpoints,
		CreatedBy:          "github.com/k4sper1love",
		License:            license,
	}

	writeJSON(w, r, http.StatusOK, envelope{"api": response})
}
