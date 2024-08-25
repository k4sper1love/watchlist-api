package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func route() *mux.Router {
	r := mux.NewRouter()

	r.Use(jwtAuth)

	r.NotFoundHandler = http.HandlerFunc(notFoundResponse)
	r.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedResponse)

	r.HandleFunc("/api/v1/healthcheck", healthcheckHandler).Methods(http.MethodGet)

	auth1 := r.PathPrefix("/api/v1").Subrouter()

	auth1.HandleFunc("/auth/register", registerHandler).Methods(http.MethodPost)
	auth1.HandleFunc("/auth/login", loginHandler).Methods(http.MethodPost)
	auth1.HandleFunc("/auth/refresh", refreshAccessTokenHandler).Methods(http.MethodPost)
	auth1.HandleFunc("/auth/logout", logoutHandler).Methods(http.MethodPost)

	user1 := r.PathPrefix("/api/v1").Subrouter()

	user1.HandleFunc("/user", getUserHandler).Methods(http.MethodGet)
	user1.HandleFunc("/user", updateUserHandler).Methods(http.MethodPut)
	user1.HandleFunc("/user", deleteUserHandler).Methods(http.MethodDelete)

	films1 := r.PathPrefix("/api/v1").Subrouter()

	films1.HandleFunc("/films", getFilmsHandler).Methods(http.MethodGet)
	films1.HandleFunc("/films", addFilmHandler).Methods(http.MethodPost)
	films1.HandleFunc("/films/{filmId:[0-9]+}", getFilmHandler).Methods(http.MethodGet)
	films1.HandleFunc("/films/{filmId:[0-9]+}", updateFilmHandler).Methods(http.MethodPut)
	films1.HandleFunc("/films/{filmId:[0-9]+}", deleteFilmHandler).Methods(http.MethodDelete)

	collections1 := r.PathPrefix("/api/v1").Subrouter()

	collections1.HandleFunc("/collections", getCollectionsHandler).Methods(http.MethodGet)
	collections1.HandleFunc("/collections", addCollectionHandler).Methods(http.MethodPost)
	collections1.HandleFunc("/collections/{collectionId:[0-9]+}", getCollectionHandler).Methods(http.MethodGet)
	collections1.HandleFunc("/collections/{collectionId:[0-9]+}", updateCollectionHandler).Methods(http.MethodPut)
	collections1.HandleFunc("/collections/{collectionId:[0-9]+}", deleteCollectionHandler).Methods(http.MethodDelete)

	collectionFilms1 := r.PathPrefix("/api/v1/collections/{collectionId:[0-9]+}").Subrouter()

	collectionFilms1.HandleFunc("/films", getCollectionFilmsHandler).Methods(http.MethodGet)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", addCollectionFilmHandler).Methods(http.MethodPost)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", getCollectionFilmHandler).Methods(http.MethodGet)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", updateCollectionFilmHandler).Methods(http.MethodPut)
	collectionFilms1.HandleFunc("/films/{filmId:[0-9]+}", deleteCollectionFilmHandler).Methods(http.MethodDelete)

	return r
}
