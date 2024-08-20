package rest

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"net/http"
)

var ErrAlreadyExists = errors.New("resource already exists")

var ErrNotFound = errors.New("resource not found")

var ErrEmptyRequest = errors.New("empty request body")

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	env := envelope{"error": message}

	err := writeJSON(w, status, env)
	if err != nil {
		log.Println(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(r, err)

	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func ConflictResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusConflict, ErrAlreadyExists.Error())
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusNotFound, ErrNotFound.Error())
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported this resource", r.Method)
	ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func mapDBError(err error) error {
	var pqErr *pq.Error

	switch {
	case errors.As(err, &pqErr) && pqErr.Code == "23505":
		return ErrAlreadyExists
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	default:
		return err
	}
}
