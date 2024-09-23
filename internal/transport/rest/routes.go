// Package rest provides HTTP handlers for managing and retrieving information related to the REST API.
//
// This package includes handlers for adding, retrieving, updating, and deleting users, films,
// collections, and collection-films, as well as for checking the health of the API.
//
// The handlers use a custom logger for logging and interact with the database and other internal
// components to perform various operations related to users, films, collections, and permissions.

package rest

import (
	"github.com/gorilla/mux"
	_ "github.com/k4sper1love/watchlist-api/api"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// route initializes the HTTP router with all API routes and handlers.
// It sets up routes for health checks, user authentication, and CRUD operations
// for users, films, collections, and collection films. It also configures middleware
// for authentication and custom handlers for not found and method not allowed errors.
//
// Returns:
//
//	*mux.Router - Configured router with all routes and middleware.
func route() *mux.Router {
	router := mux.NewRouter() //Register a new router

	router.Use(requireAuth) // Apply JWT authentication middleware

	// Handle 404 Not Found
	router.NotFoundHandler = http.HandlerFunc(notFoundResponse)

	// Handle 405 Method Not Allowed
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedResponse)

	// Default endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api", http.StatusSeeOther)
	})

	// Info API endpoint
	router.HandleFunc("/api", infoHandler).Methods(http.MethodGet)

	// Health check endpoint
	router.HandleFunc("/api/v1/healthcheck", healthcheckHandler).Methods(http.MethodGet)

	// Swagger documentation UI endpoint
	router.Handle("/swagger/{rest:.*}", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Swagger UI will use the Swagger JSON file
	))

	// Subrouter for authentication routes
	auth1 := router.PathPrefix("/api/v1").Subrouter()
	auth1.HandleFunc("/auth/register", registerHandler).Methods(http.MethodPost)
	auth1.HandleFunc("/auth/login", loginHandler).Methods(http.MethodPost)
	auth1.HandleFunc("/auth/refresh", refreshAccessTokenHandler).Methods(http.MethodPost)
	auth1.HandleFunc("/auth/logout", logoutHandler).Methods(http.MethodPost)

	// Subrouter for user routes
	user1 := router.PathPrefix("/api/v1").Subrouter()
	user1.HandleFunc("/user", getUserHandler).Methods(http.MethodGet)
	user1.HandleFunc("/user", updateUserHandler).Methods(http.MethodPut)
	user1.HandleFunc("/user", deleteUserHandler).Methods(http.MethodDelete)

	// Subrouter for film routes
	films1 := router.PathPrefix("/api/v1").Subrouter()
	films1.HandleFunc("/films", getFilmsHandler).Methods(http.MethodGet)
	films1.HandleFunc("/films", requirePermissions("film", "create", addFilmHandler)).Methods(http.MethodPost)
	films1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("film", "read", getFilmHandler)).Methods(http.MethodGet)
	films1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("film", "update", updateFilmHandler)).Methods(http.MethodPut)
	films1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("film", "delete", deleteFilmHandler)).Methods(http.MethodDelete)

	// Subrouter for collection routes
	collections1 := router.PathPrefix("/api/v1").Subrouter()
	collections1.HandleFunc("/collections", getCollectionsHandler).Methods(http.MethodGet)
	collections1.HandleFunc("/collections", requirePermissions("collection", "create", addCollectionHandler)).Methods(http.MethodPost)
	collections1.HandleFunc("/collections/{collectionId:[0-9]+}", requirePermissions("collection", "read", getCollectionHandler)).Methods(http.MethodGet)
	collections1.HandleFunc("/collections/{collectionId:[0-9]+}", requirePermissions("collection", "update", updateCollectionHandler)).Methods(http.MethodPut)
	collections1.HandleFunc("/collections/{collectionId:[0-9]+}", requirePermissions("collection", "delete", deleteCollectionHandler)).Methods(http.MethodDelete)

	// Subrouter for collection films routes
	collectionFilms1 := router.PathPrefix("/api/v1/collections/{collectionId:[0-9]+}").Subrouter()
	collectionFilms1.HandleFunc("/films", requirePermissions("collectionFilm", "read", getCollectionFilmsHandler)).Methods(http.MethodGet)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("collectionFilm", "create", addCollectionFilmHandler)).Methods(http.MethodPost)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("collectionFilm", "read", getCollectionFilmHandler)).Methods(http.MethodGet)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("collectionFilm", "update", updateCollectionFilmHandler)).Methods(http.MethodPut)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("collectionFilm", "delete", deleteCollectionFilmHandler)).Methods(http.MethodDelete)

	return router
}
