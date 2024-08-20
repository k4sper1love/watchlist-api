package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func route() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", defaultHandler).Methods(http.MethodGet)
	r.NotFoundHandler = http.HandlerFunc(NotFoundResponse)
	r.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedResponse)

	users := r.PathPrefix("/users").Subrouter()

	users.HandleFunc("", getUsersHandler).Methods(http.MethodGet)
	users.HandleFunc("", addUserHandler).Methods(http.MethodPost)

	users.HandleFunc("/{id:[0-9]+}", getUserHandler).Methods(http.MethodGet)
	users.HandleFunc("/{id:[0-9]+}", updateUserHandler).Methods(http.MethodPut)
	users.HandleFunc("/{id:[0-9]+}", deleteUserHandler).Methods(http.MethodDelete)

	return r
}
