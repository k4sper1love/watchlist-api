package rest

import (
	"database/sql"
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/validator"
	"net/http"
)

// GetUser godoc
// @Summary Get user account
// @Description Get information about user by ID using an authentication token.
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} swagger.UserResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /user [get]
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

// UpdateUser godoc
// @Summary Update user account
// @Description Update information about user by ID using an authentication token.
// @Tags user
// @Accept json
// @Produce json
// @Param username body swagger.UpdateUserRequest true "New information about the user"
// @Success 200 {object} swagger.UserResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /user [put]
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

// DeleteUser godoc
// @Summary Delete user account
// @Description Delete user by ID using an authentication token.
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} swagger.MessageResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /user [delete]
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
