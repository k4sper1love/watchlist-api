package rest

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Predefined error messages
var (
	errAlreadyExists       = errors.New("resource already exists")
	errNotFound            = errors.New("resource not found")
	errEmptyRequest        = errors.New("empty request body")
	errForeignKeyViolation = errors.New("attempted to reference a non-existent record")
	errInvalidRefreshToken = errors.New("invalid or revoked refresh token")
)

// errorResponse sends a JSON response with an error message and status code.
func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	writeJSON(w, r, status, envelope{"error": message})
}

// serverErrorResponse handles internal server errors, logging the error and sending a generic message to the client.
func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
	sl.PrintEndpointError("server error", err, r)
}

// badRequestResponse handles client errors where the request is invalid.
func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
	sl.PrintEndpointWarn("bad request", err, r)
}

// uniqueConflictResponse handles errors where a resource already exists in the database.
func uniqueConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusConflict, err.Error())
	sl.PrintEndpointWarn("unique conflict", err, r)
}

// editConflictResponse handles conflicts when updating a record.
func editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "record update failed due to a conflict. Please try again"
	errorResponse(w, r, http.StatusConflict, message)
	sl.PrintEndpointWarn("edit conflict", nil, r)
}

// notFoundResponse handles cases where a requested resource is not found.
func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusNotFound, errNotFound.Error())
	sl.PrintEndpointWarn("not found", nil, r)
}

// methodNotAllowedResponse handles requests with unsupported HTTP methods.
func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
	sl.PrintEndpointWarn(message, nil, r)
}

// invalidAuthTokenResponse handles cases where the authentication token is invalid or missing.
func invalidAuthTokenResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid or missing authentication token"
	errorResponse(w, r, http.StatusUnauthorized, message)
	sl.PrintEndpointWarn(message, nil, r)
}

// incorrectPasswordResponse handles authentication failures due to incorrect passwords.
func incorrectPasswordResponse(w http.ResponseWriter, r *http.Request) {
	message := "incorrect password"
	errorResponse(w, r, http.StatusUnauthorized, message)
	sl.PrintEndpointWarn(message, nil, r)
}

// failedValidationResponse handles cases where input validation fails.
func failedValidationResponse(w http.ResponseWriter, r *http.Request, errs map[string]string) {
	errorResponse(w, r, http.StatusUnprocessableEntity, errs)
	sl.PrintEndpointWarn("validation errors", errs, r)
}

// forbiddenResponse handles requests where the user does not have the required permissions.
func forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	message := "you don't have enough permissions to perform this action"
	errorResponse(w, r, http.StatusForbidden, message)
	sl.PrintEndpointWarn(message, nil, r)
}

// handleDBError processes database errors and maps them to appropriate HTTP responses.
func handleDBError(w http.ResponseWriter, r *http.Request, err error) {
	var pqErr *pq.Error

	switch {
	case errors.As(err, &pqErr):
		if pqErr.Code == "23505" {
			uniqueConflictResponse(w, r, errAlreadyExists)
			return
		}
		if pqErr.Code == "23503" {
			uniqueConflictResponse(w, r, errForeignKeyViolation)
			return
		}
		serverErrorResponse(w, r, err)
		return
	case errors.Is(err, sql.ErrNoRows):
		notFoundResponse(w, r)
		return
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		incorrectPasswordResponse(w, r)
		return
	default:
		serverErrorResponse(w, r, err)
		return
	}
}
