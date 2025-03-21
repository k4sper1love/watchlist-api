package rest

import (
	"github.com/gorilla/mux"
	_ "github.com/k4sper1love/watchlist-api/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// route initializes the HTTP router with all API routes and handlers.
func route() *mux.Router {
	// Register a new router
	router := mux.NewRouter()

	// Apply middlewares
	router.Use(logAndRecordMetrics)
	router.Use(authenticate)

	// API Icon Endpoint
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	router.HandleFunc("/upload", uploadImageHandler).Methods(http.MethodPost)

	router.HandleFunc("/images/{filename}", getImageHandler).Methods(http.MethodGet)

	// Handle 404 Not Found
	router.NotFoundHandler = http.HandlerFunc(notFoundResponse)

	// Handle 405 Method Not Allowed
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedResponse)

	// Default Endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api", http.StatusSeeOther)
	})

	// Info API Endpoint
	router.HandleFunc("/api", infoHandler).Methods(http.MethodGet)

	// Health check Endpoint
	router.HandleFunc("/api/v1/healthcheck", healthcheckHandler).Methods(http.MethodGet)

	// Swagger documentation UI Endpoint
	router.HandleFunc("/swagger/{rest:.*}", swaggerHandler)

	// Prometheus Metrics Endpoint
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// Set up routes
	setupAuthRoutes(router)
	setupUserRoutes(router)
	setupFilmRoutes(router)
	setupCollectionRoutes(router)
	setupCollectionFilmRoutes(router)

	return router
}

func setupAuthRoutes(router *mux.Router) {
	auth := router.PathPrefix("/api/v1/auth").Subrouter()
	auth.HandleFunc("/register", registerWithCredentialsHandler).Methods(http.MethodPost)
	auth.HandleFunc("/register/telegram", verificate(registerByTelegramHandler)).Methods(http.MethodPost)
	auth.HandleFunc("/login", loginWithCredentialsHandler).Methods(http.MethodPost)
	auth.HandleFunc("/login/telegram", verificate(loginByTelegramHandler)).Methods(http.MethodPost)
	auth.HandleFunc("/refresh", refreshAccessTokenHandler).Methods(http.MethodPost)
	auth.HandleFunc("/logout", logoutHandler).Methods(http.MethodPost)
	auth.HandleFunc("/check", checkTokenHandler).Methods(http.MethodGet)
}

func setupUserRoutes(router *mux.Router) {
	user := router.PathPrefix("/api/v1").Subrouter()
	user.HandleFunc("/user", getUserHandler).Methods(http.MethodGet)
	user.HandleFunc("/user", updateUserHandler).Methods(http.MethodPut)
	user.HandleFunc("/user", deleteUserHandler).Methods(http.MethodDelete)
}

func setupFilmRoutes(router *mux.Router) {
	films := router.PathPrefix("/api/v1/films").Subrouter()
	films.HandleFunc("", getFilmsHandler).Methods(http.MethodGet)
	films.HandleFunc("", requirePermissions("film", "create", addFilmHandler)).Methods(http.MethodPost)
	films.HandleFunc("/{filmID:[0-9]+}", requirePermissions("film", "read", getFilmHandler)).Methods(http.MethodGet)
	films.HandleFunc("/{filmID:[0-9]+}", requirePermissions("film", "update", updateFilmHandler)).Methods(http.MethodPut)
	films.HandleFunc("/{filmID:[0-9]+}", requirePermissions("film", "delete", deleteFilmHandler)).Methods(http.MethodDelete)
}

func setupCollectionRoutes(router *mux.Router) {
	collections := router.PathPrefix("/api/v1/collections").Subrouter()
	collections.HandleFunc("", getCollectionsHandler).Methods(http.MethodGet)
	collections.HandleFunc("", requirePermissions("collection", "create", addCollectionHandler)).Methods(http.MethodPost)
	collections.HandleFunc("/{collectionID:[0-9]+}", requirePermissions("collection", "read", getCollectionHandler)).Methods(http.MethodGet)
	collections.HandleFunc("/{collectionID:[0-9]+}", requirePermissions("collection", "update", updateCollectionHandler)).Methods(http.MethodPut)
	collections.HandleFunc("/{collectionID:[0-9]+}", requirePermissions("collection", "delete", deleteCollectionHandler)).Methods(http.MethodDelete)
}

func setupCollectionFilmRoutes(router *mux.Router) {
	collectionFilms := router.PathPrefix("/api/v1/collections/{collectionID:[0-9]+}/films").Subrouter()
	collectionFilms.HandleFunc("", requirePermissions("collectionFilm", "read", getCollectionFilmsHandler)).Methods(http.MethodGet)
	collectionFilms.HandleFunc("", requirePermissions("collectionFilm", "create", addNewCollectionFilmHandler)).Methods(http.MethodPost)
	collectionFilms.HandleFunc("/{filmID:[0-9]+}", requirePermissions("collectionFilm", "add", addCollectionFilmHandler)).Methods(http.MethodPost)
	collectionFilms.HandleFunc("/{filmID:[0-9]+}", requirePermissions("collectionFilm", "read", getCollectionFilmHandler)).Methods(http.MethodGet)
	collectionFilms.HandleFunc("/{filmID:[0-9]+}", requirePermissions("collectionFilm", "delete", deleteCollectionFilmHandler)).Methods(http.MethodDelete)
}
