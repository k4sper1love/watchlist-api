package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func route() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", defaultHandler).Methods(http.MethodGet)
	r.NotFoundHandler = http.HandlerFunc(notFoundResponse)
	r.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedResponse)

	users := r.PathPrefix("/users").Subrouter()

	users.HandleFunc("", getUsersHandler).Methods(http.MethodGet)
	users.HandleFunc("", addUserHandler).Methods(http.MethodPost)
	users.HandleFunc("/{userId:[0-9]+}", getUserHandler).Methods(http.MethodGet)
	users.HandleFunc("/{userId:[0-9]+}", updateUserHandler).Methods(http.MethodPut)
	users.HandleFunc("/{userId:[0-9]+}", deleteUserHandler).Methods(http.MethodDelete)

	films := r.PathPrefix("/films").Subrouter()

	films.HandleFunc("", getFilmsHandler).Methods(http.MethodGet)
	films.HandleFunc("", addFilmHandler).Methods(http.MethodPost)
	films.HandleFunc("/{filmId:[0-9]+}", getFilmHandler).Methods(http.MethodGet)
	films.HandleFunc("/{filmId:[0-9]+}", updateFilmHandler).Methods(http.MethodPut)
	films.HandleFunc("/{filmId:[0-9]+}", deleteFilmHandler).Methods(http.MethodDelete)

	collections := r.PathPrefix("/users/{userId:[0-9]+}/collections").Subrouter()

	collections.HandleFunc("", getCollectionsHandler).Methods(http.MethodGet)
	collections.HandleFunc("", addCollectionHandler).Methods(http.MethodPost)
	collections.HandleFunc("/{collectionId:[0-9]+}", getCollectionHandler).Methods(http.MethodGet)
	collections.HandleFunc("/{collectionId:[0-9]+}", updateCollectionHandler).Methods(http.MethodPut)
	collections.HandleFunc("/{collectionId:[0-9]+}", deleteCollectionHandler).Methods(http.MethodDelete)

	collectionFilms := collections.PathPrefix("/{collectionId:[0-9]+}/films").Subrouter()

	collectionFilms.HandleFunc("", getCollectionFilmsHandler).Methods(http.MethodGet)
	collectionFilms.HandleFunc("/{filmId:[0-9]+}", addCollectionFilmHandler).Methods(http.MethodPost)
	collectionFilms.HandleFunc("/{filmId:[0-9]+}", getCollectionFilmHandler).Methods(http.MethodGet)
	collectionFilms.HandleFunc("/{filmId:[0-9]+}", updateCollectionFilmHandler).Methods(http.MethodPut)
	collectionFilms.HandleFunc("/{filmId:[0-9]+}", deleteCollectionFilmHandler).Methods(http.MethodDelete)

	viewedFilms := r.PathPrefix("/users/{userId:[0-9]+}/viewed").Subrouter()

	viewedFilms.HandleFunc("", getViewedFilmsHandler).Methods(http.MethodGet)
	viewedFilms.HandleFunc("/{filmId:[0-9]+}", addViewedFilmHandler).Methods(http.MethodPost)
	viewedFilms.HandleFunc("/{filmId:[0-9]+}", getViewedFilmHandler).Methods(http.MethodGet)
	viewedFilms.HandleFunc("/{filmId:[0-9]+}", updateViewedFilmHandler).Methods(http.MethodPut)
	viewedFilms.HandleFunc("/{filmId:[0-9]+}", deleteViewedFilmHandler).Methods(http.MethodDelete)

	return r
}
