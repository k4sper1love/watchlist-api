package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/models"
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
	var user models.User

	if err := parseRequestBody(r, &user); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if errs := validator.ValidateStruct(&user); errs != nil {
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

	if err = postgres.AddUserPermissions(user.Id, permissionCodes...); err != nil {
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
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := parseRequestBody(r, &input); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if errs := validator.ValidateStruct(&input); errs != nil {
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
	refreshToken := parseTokenFromHeader(r)
	if refreshToken == "" {
		invalidAuthTokenResponse(w, r)
		return
	}

	// Revoke the refresh token.
	if err := logout(refreshToken); err != nil {
		invalidAuthTokenResponse(w, r)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "token revoked"})
}
