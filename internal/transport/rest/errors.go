package rest

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var errAlreadyExists = errors.New("resource already exists")

var errNotFound = errors.New("resource not found")

var errEmptyRequest = errors.New("empty request body")

var errForeignKeyViolation = errors.New("attempted to reference a non-existent record")

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	env := envelope{"error": message}

	writeJSON(w, r, status, env)
}

func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(r, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusConflict, err.Error())
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusNotFound, errNotFound.Error())
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func invalidAuthTokenResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid or missing authentication token"
	errorResponse(w, r, http.StatusUnauthorized, message)
}

func invalidLoginResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid login credentials"
	errorResponse(w, r, http.StatusUnauthorized, message)
}

func handleDBError(w http.ResponseWriter, r *http.Request, err error) {
	var pqErr *pq.Error

	switch {
	case errors.As(err, &pqErr):
		if pqErr.Code == "23505" {
			conflictResponse(w, r, errAlreadyExists)
			return
		}
		if pqErr.Code == "23503" {
			conflictResponse(w, r, errForeignKeyViolation)
			return
		}
		serverErrorResponse(w, r, err)
		return
	case errors.Is(err, sql.ErrNoRows):
		notFoundResponse(w, r)
		return
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		invalidLoginResponse(w, r)
		return
	default:
		serverErrorResponse(w, r, err)
		return
	}
}
