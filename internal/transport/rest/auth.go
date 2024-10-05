package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// register creates a new user and generates authentication tokens.
func register(user *models.User) (*models.AuthResponse, error) {
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

	res := &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	return res, nil
}

// login authenticates a user by email and password, and generates authentication tokens.
func login(email, password string) (*models.AuthResponse, error) {
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

func checkToken(token string) error {
	if claims := parseTokenClaims(token); claims == nil {
		return errInvalidRefreshToken
	}
	return nil
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
func createAuthResponse(user *models.User, accessToken, refreshToken string) *models.AuthResponse {
	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
