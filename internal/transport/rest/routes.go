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

	films := r.PathPrefix("/films").Subrouter()

	films.HandleFunc("", getFilmsHandler).Methods(http.MethodGet)
	films.HandleFunc("", addFilmHandler).Methods(http.MethodPost)

	film := films.PathPrefix("/{filmId:[0-9]+}").Subrouter()

	film.HandleFunc("", getFilmHandler).Methods(http.MethodGet)
	film.HandleFunc("", updateFilmHandler).Methods(http.MethodPut)
	film.HandleFunc("", deleteFilmHandler).Methods(http.MethodDelete)

	users := r.PathPrefix("/users").Subrouter()

	users.HandleFunc("", getUsersHandler).Methods(http.MethodGet)
	users.HandleFunc("", addUserHandler).Methods(http.MethodPost)

	user := users.PathPrefix("/{userId:[0-9]+}").Subrouter()

	user.HandleFunc("", getUserHandler).Methods(http.MethodGet)
	user.HandleFunc("", updateUserHandler).Methods(http.MethodPut)
	user.HandleFunc("", deleteUserHandler).Methods(http.MethodDelete)

	userFilms := user.PathPrefix("/films").Subrouter()

	userFilms.HandleFunc("", getUserFilmsHandler).Methods(http.MethodGet)
	userFilms.HandleFunc("", addUserFilm).Methods(http.MethodPost)
	//userFilms.HandleFunc("/new", defaultHandler).Methods(http.MethodPost) - no ideas now

	userFilm := userFilms.PathPrefix("/{filmId:[0-9]+}").Subrouter()

	userFilm.HandleFunc("", getUserFilmHandler).Methods(http.MethodGet)
	userFilm.HandleFunc("", updateUserFilmHandler).Methods(http.MethodPut)
	userFilm.HandleFunc("", deleteUserFilmHandler).Methods(http.MethodDelete)

	//viewedFilms := r.PathPrefix("/users/{userId:[0-9]+/viewed}").Subrouter()
	//
	//viewedFilms.HandleFunc("", defaultHandler).Methods(http.MethodGet)
	//viewedFilms.HandleFunc("", defaultHandler).Methods(http.MethodPost)
	//
	//viewedFilm.HandleFunc("/{id:[0-9]+}", defaultHandler).Methods(http.MethodGet)
	//viewedFilm.HandleFunc("/{id:[0-9]+}", defaultHandler).Methods(http.MethodPut)
	//viewedFilm.HandleFunc("/{id:[0-9]+}", defaultHandler).Methods(http.MethodDelete)

	return r
}
