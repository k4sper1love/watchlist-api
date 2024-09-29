package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// authResponse represents the response returned after authentication operations.
type authResponse struct {
	*models.User
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// register creates a new user and generates authentication tokens.
func register(user *models.User) (*authResponse, error) {
	if err := hashPassword(user); err != nil {
		return nil, err
	}

	if err := postgres.AddUser(user); err != nil {
		return nil, err
	}

	user.Password = "" // Clear the password before returning.

	accessToken, err := generateAccessToken(user.Id)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateAndSaveRefreshToken(user.Id)
	if err != nil {
		return nil, err
	}

	res := &authResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	return res, nil
}

// login authenticates a user by email and password, and generates authentication tokens.
func login(email, password string) (*authResponse, error) {
	// Retrieve the user from the database by email.
	user, err := postgres.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := comparePasswords(user.Password, password); err != nil {
		return nil, err
	}

	user.Password = "" // Clear the password before returning.

	accessToken, err := generateAccessToken(user.Id)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateAndSaveRefreshToken(user.Id)
	if err != nil {
		return nil, err
	}

	return createAuthResponse(user, accessToken, refreshToken), nil
}

// refreshAccessToken generates a new access token using a valid refresh token.
func refreshAccessToken(refreshToken string) (string, error) {
	if claims := parseTokenClaims(refreshToken); claims == nil {
		return "", errInvalidRefreshToken
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
	if claims := parseTokenClaims(refreshToken); claims == nil {
		return errInvalidRefreshToken
	}

	if isRevoked, err := postgres.IsRefreshTokenRevoked(refreshToken); err != nil || isRevoked {
		return errInvalidRefreshToken
	}

	return postgres.RevokeRefreshToken(refreshToken)
}

// hashPassword hashes the user's password using bcrypt.
func hashPassword(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// comparePasswords compares a hashed password with a plaintext password.
func comparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// createAuthResponse constructs and returns an authResponse object.
func createAuthResponse(user *models.User, accessToken, refreshToken string) *authResponse {
	return &authResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
