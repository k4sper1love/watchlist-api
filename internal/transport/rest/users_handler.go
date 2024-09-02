package rest

import (
	"database/sql"
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/validator"
	"net/http"
)

// getUserHandler retrieves the user details for the authenticated user.
//
// Returns a JSON response with the user details or an error if retrieval fails.
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Retrieve the user ID from the request context.
	userId := r.Context().Value("userId").(int)

	user, err := postgres.GetUserById(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}
	// Remove the password from the user details before sending the response.
	user.Password = ""

	writeJSON(w, r, http.StatusOK, envelope{"user": user})
}

// updateUserHandler updates the details of the authenticated user.
//
// Returns a JSON response with the updated user details or an error if the update fails.
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Retrieve the user ID from the request context.
	userId := r.Context().Value("userId").(int)

	user, err := postgres.GetUserById(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = parseRequestBody(r, user)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// Validate the updated user details.
	errs := validator.ValidateStruct(user)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}
	// Set the user ID and clear the password field.
	user.Id = userId
	user.Password = ""

	err = postgres.UpdateUser(user)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			editConflictResponse(w, r)
		default:
			handleDBError(w, r, err)
		}
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"user": user})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Retrieve the user ID from the request context.
	userId := r.Context().Value("userId").(int)

	// Verify that the user exists in the database.
	_, err := postgres.GetUserById(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = postgres.DeleteUser(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	// Confirm successful deletion with a JSON response.
	writeJSON(w, r, http.StatusOK, envelope{"message": "user deleted"})
}
