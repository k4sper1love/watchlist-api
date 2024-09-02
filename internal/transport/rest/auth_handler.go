package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/validator"
	"net/http"
)

// registerHandler registers a new user.
//
// Returns a JSON response with the created user or an error if registration fails.
func registerHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	var user models.User
	err := parseRequestBody(r, &user)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	errs := validator.ValidateStruct(&user)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	// Register the user in the system.
	resp, err := register(&user)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	// Assign default permissions to the user.
	permissionCodes := []string{"film:create", "collection:create"}
	err = postgres.AddUserPermissions(user.Id, permissionCodes...)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"user": resp})
}

// loginHandler logs in a user.
//
// Returns a JSON response with user information or an error if login fails.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	var input struct {
		Email    string `json:"email" validate:"required,email"` // Email of the user attempting to log in.
		Password string `json:"password" validate:"required"`    // Password of the user attempting to log in.
	}

	err := parseRequestBody(r, &input)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	errs := validator.ValidateStruct(&input)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	// Authenticate the user.
	resp, err := login(input.Email, input.Password)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"user": resp})
}

// refreshAccessTokenHandler refreshes the access token.
//
// Returns a JSON response with a new access token or an error if refresh fails.
func refreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Extract the refresh token from the request header.
	refreshToken := parseTokenFromHeader(r)
	if refreshToken == "" {
		invalidAuthTokenResponse(w, r)
		return
	}

	// Refresh the access token.
	newAccessToken, err := refreshAccessToken(refreshToken)
	if err != nil {
		invalidAuthTokenResponse(w, r)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"access_token": newAccessToken})
}

// logoutHandler logs out a user by revoking the refresh token.
//
// Returns a JSON response confirming token revocation or an error if logout fails.
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Extract the refresh token from the request header.
	refreshToken := parseTokenFromHeader(r)
	if refreshToken == "" {
		invalidAuthTokenResponse(w, r)
		return
	}

	// Revoke the refresh token.
	err := logout(refreshToken)
	if err != nil {
		invalidAuthTokenResponse(w, r)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "token revoked"})
}
