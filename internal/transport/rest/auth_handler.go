package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/validator"
	"net/http"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user using a username, email and password. Returns user information and tokens.
// @Description Basic permissions are available to you: creating films and collections.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body swagger.RegisterRequest true "Information about the new user "
// @Success 201 {object} swagger.AuthResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Router /auth/register [post]
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

// Login godoc
// @Summary Log in to your account
// @Description Log in to your account using your email and password. Returns tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body swagger.LoginRequest true "Login information"
// @Success 200 {object} swagger.AuthResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Router /auth/login [post]
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

// Refresh godoc
// @Summary Refresh access token
// @Description Refresh your access token using the refresh token in the Authorization header.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} swagger.AccessTokenResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /auth/refresh [post]
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

// Logout godoc
// @Summary Log out of your account
// @Description Log out of your account using your refresh token in the Authorization header.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} swagger.MessageResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /auth/logout [post]
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
