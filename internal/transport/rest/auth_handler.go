package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-api/pkg/validator"
	"net/http"
)

// RegisterWithCredentials godoc
// @Summary Register a new user with credentials
// @Description Register a new user using a username and password. Returns user information and tokens.
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
func registerWithCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	if err := parseRequestBody(r, &credentials); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if errs := validator.ValidateStruct(&credentials); errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	// Register the user in the system.
	user, err := registerWithCredentials(&credentials)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	// Assign default permissions to the user.
	permissionCodes := []string{"film:create", "collection:create"}

	if err = postgres.AddUserPermissions(user.ID, permissionCodes...); err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"user": user})
}

// RegisterByTelegram godoc
// @Summary Register a new user by Telegram
// @Description Register a new user using verification token from header. Returns user information and tokens.
// @Description Basic permissions are available to you: creating films and collections.
// @Tags auth
// @Accept json
// @Produce json
// @Param Verification header string true "Verification token from Telegram"
// @Success 201 {object} swagger.AuthResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Router /auth/register/telegram [post]
func registerByTelegramHandler(w http.ResponseWriter, r *http.Request) {
	telegramID := r.Context().Value("telegramID").(int)

	// Register the user in the system.
	user, err := registerByTelegram(telegramID)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	// Assign default permissions to the user.
	permissionCodes := []string{"film:create", "collection:create"}

	if err = postgres.AddUserPermissions(user.ID, permissionCodes...); err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"user": user})
}

// LoginWithCredentials godoc
// @Summary Log in to your account with credentials
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
func loginWithCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	if err := parseRequestBody(r, &credentials); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if errs := validator.ValidateStruct(&credentials); errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	// Authenticate the user.
	resp, err := loginWithCredentials(credentials.Username, credentials.Password)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"user": resp})
}

// LoginByTelegram godoc
// @Summary Log in to your account by Telegram
// @Description Log in to your account using verification token from header. Returns tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param Verification header string true "Verification token from Telegram"
// @Success 200 {object} swagger.AuthResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Router /auth/login/telegram [post]
func loginByTelegramHandler(w http.ResponseWriter, r *http.Request) {
	telegramID := r.Context().Value("telegramID").(int)

	// Authenticate the user.
	resp, err := loginByTelegram(telegramID)
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

// CheckToken godoc
// @Summary Check validity of token
// @Description Checks if the token provided in the Authorization header is still valid.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} swagger.MessageResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /auth/check-token [get]
func checkTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := parseTokenFromHeader(r)
	if token == "" {
		invalidAuthTokenResponse(w, r)
		return
	}

	if err := checkToken(token); err != nil {
		invalidAuthTokenResponse(w, r)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "token is valid"})
}
