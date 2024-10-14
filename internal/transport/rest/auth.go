package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// registerWithCredentials creates a new user using provided credentials and generates authentication tokens.
func registerWithCredentials(credentials *models.Credentials) (*models.AuthResponse, error) {
	if err := hashPassword(credentials); err != nil {
		return nil, err
	}

	user, err := postgres.AddUserWithCredentials(credentials)
	if err != nil {
		return nil, err
	}

	user.Password = "" // Clear the password before returning.

	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateAndSaveRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return createAuthResponse(user, accessToken, refreshToken), nil
}

// registerByTelegram creates a new user with the provided Telegram ID and generates authentication tokens.
func registerByTelegram(telegramID int) (*models.AuthResponse, error) {
	credentials := &models.Credentials{
		TelegramID: telegramID,
		Username:   generateUniqueUsername(4, telegramID),
	}

	user, err := postgres.AddUserByTelegramID(credentials)
	if err != nil {
		return nil, err
	}

	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateAndSaveRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return createAuthResponse(user, accessToken, refreshToken), nil
}

// loginWithCredentials authenticates a user by their username and password, generating authentication tokens upon success.
func loginWithCredentials(username, password string) (*models.AuthResponse, error) {
	// Retrieve the user from the database by email.
	user, err := postgres.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if user.Password == "" {
		return nil, errRequiredPassword
	}

	if err := comparePasswords(user.Password, password); err != nil {
		return nil, err
	}

	user.Password = "" // Clear the password before returning.

	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateAndSaveRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return createAuthResponse(user, accessToken, refreshToken), nil
}

// loginByTelegram authenticates a user using their Telegram ID and generates authentication tokens.
func loginByTelegram(telegramID int) (*models.AuthResponse, error) {
	// Retrieve the user from the database by email.
	user, err := postgres.GetUserByTelegramID(telegramID)
	if err != nil {
		return nil, err
	}

	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateAndSaveRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return createAuthResponse(user, accessToken, refreshToken), nil
}

// refreshAccessToken generates a new access token using a valid refresh token.
func refreshAccessToken(refreshToken string) (string, error) {
	if err := checkToken(refreshToken); err != nil {
		return "", err
	}

	if isRevoked, err := postgres.IsRefreshTokenRevoked(refreshToken); err != nil || isRevoked {
		return "", errInvalidRefreshToken
	}

	userId, err := postgres.GetIdFromRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return generateAccessToken(userId)
}

// logout invalidates the given refresh token.
func logout(refreshToken string) error {
	if err := checkToken(refreshToken); err != nil {
		return err
	}

	if isRevoked, err := postgres.IsRefreshTokenRevoked(refreshToken); err != nil || isRevoked {
		return errInvalidRefreshToken
	}

	return postgres.RevokeRefreshToken(refreshToken)
}

// checkToken verifies the validity of the provided authentication token.
func checkToken(token string) error {
	claims, err := parseTokenClaims(token, config.JWTSecret)
	if err != nil {
		return err
	}

	if claims == nil {
		return errInvalidRefreshToken
	}
	return nil
}

// hashPassword hashes the user's password using bcrypt.
func hashPassword(credentials *models.Credentials) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	credentials.Password = string(hashedPassword)
	return nil
}

// comparePasswords compares a hashed password with a plaintext password.
func comparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// createAuthResponse constructs and returns an authResponse object.
func createAuthResponse(user *models.User, accessToken, refreshToken string) *models.AuthResponse {
	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
