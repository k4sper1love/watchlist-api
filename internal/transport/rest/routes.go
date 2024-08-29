package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func route() *mux.Router {
	router := mux.NewRouter()

	router.Use(jwtAuth)

	router.NotFoundHandler = http.HandlerFunc(notFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedResponse)

	router.HandleFunc("/api/v1/healthcheck", healthcheckHandler).Methods(http.MethodGet)

	auth1 := router.PathPrefix("/api/v1").Subrouter()

	auth1.HandleFunc("/auth/register", registerHandler).Methods(http.MethodPost)
	auth1.HandleFunc("/auth/login", loginHandler).Methods(http.MethodPost)
	auth1.HandleFunc("/auth/refresh", refreshAccessTokenHandler).Methods(http.MethodPost)
	auth1.HandleFunc("/auth/logout", logoutHandler).Methods(http.MethodPost)

	user1 := router.PathPrefix("/api/v1").Subrouter()

	user1.HandleFunc("/user", getUserHandler).Methods(http.MethodGet)
	user1.HandleFunc("/user", updateUserHandler).Methods(http.MethodPut)
	user1.HandleFunc("/user", deleteUserHandler).Methods(http.MethodDelete)

	films1 := router.PathPrefix("/api/v1").Subrouter()

	films1.HandleFunc("/films", getFilmsHandler).Methods(http.MethodGet)
	films1.HandleFunc("/films", requirePermissions("film", "create", addFilmHandler)).Methods(http.MethodPost)
	films1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("film", "read", getFilmHandler)).Methods(http.MethodGet)
	films1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("film", "update", updateFilmHandler)).Methods(http.MethodPut)
	films1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("film", "delete", deleteFilmHandler)).Methods(http.MethodDelete)

	collections1 := router.PathPrefix("/api/v1").Subrouter()

	collections1.HandleFunc("/collections", getCollectionsHandler).Methods(http.MethodGet)
	collections1.HandleFunc("/collections", requirePermissions("collection", "create", addCollectionHandler)).Methods(http.MethodPost)
	collections1.HandleFunc("/collections/{collectionId:[0-9]+}", requirePermissions("collection", "read", getCollectionHandler)).Methods(http.MethodGet)
	collections1.HandleFunc("/collections/{collectionId:[0-9]+}", requirePermissions("collection", "update", updateCollectionHandler)).Methods(http.MethodPut)
	collections1.HandleFunc("/collections/{collectionId:[0-9]+}", requirePermissions("collection", "delete", deleteCollectionHandler)).Methods(http.MethodDelete)

	collectionFilms1 := router.PathPrefix("/api/v1/collections/{collectionId:[0-9]+}").Subrouter()

	collectionFilms1.HandleFunc("/films", requirePermissions("collectionFilm", "read", getCollectionFilmsHandler)).Methods(http.MethodGet)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("collectionFilm", "create", addCollectionFilmHandler)).Methods(http.MethodPost)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("collectionFilm", "read", getCollectionFilmHandler)).Methods(http.MethodGet)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("collectionFilm", "update", updateCollectionFilmHandler)).Methods(http.MethodPut)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", requirePermissions("collectionFilm", "delete", deleteCollectionFilmHandler)).Methods(http.MethodDelete)

	return router
}
