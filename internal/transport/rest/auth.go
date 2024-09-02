package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// authResponse represents the response returned after authentication operations.
type authResponse struct {
	*models.User
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// register creates a new user and generates authentication tokens.
//
// Returns an authResponse containing user details and tokens, or an error if registration fails.
func register(user *models.User) (*authResponse, error) {
	// Hash the user's password using bcrypt with default cost.
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword) // Store the hashed password in the user object.

	err := postgres.AddUser(user)
	if err != nil {
		return nil, err
	}
	user.Password = "" // Clear the password from the user object before returning.

	// Generate an access token for the user.
	accessToken, err := generateAccessToken(user.Id)
	if err != nil {
		return nil, err
	}

	// Generate and save a refresh token for the user.
	refreshToken, err := generateAndSaveRefreshToken(user.Id)
	if err != nil {
		return nil, err
	}

	// Create an authResponse with user details and tokens.
	res := &authResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	return res, nil
}

// login authenticates a user by email and password, and generates authentication tokens.
//
// Returns an authResponse containing user details and tokens, or an error if authentication fails.
func login(email, password string) (*authResponse, error) {
	// Retrieve the user from the database by email.
	user, err := postgres.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	user.Password = "" // Clear the password from the user object before returning.

	// Generate an access token for the user.
	accessToken, err := generateAccessToken(user.Id)
	if err != nil {
		return nil, err
	}

	// Generate and save a refresh token for the user.
	refreshToken, err := generateAndSaveRefreshToken(user.Id)
	if err != nil {
		return nil, err
	}

	// Create an authResponse with user details and tokens.
	res := &authResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	return res, nil
}

// refreshAccessToken generates a new access token using a valid refresh token.
//
// Returns the new access token or an error if the refresh token is invalid or revoked.
func refreshAccessToken(refreshToken string) (string, error) {
	// Parse the claims from the refresh token.
	claims := parseTokenClaims(refreshToken)
	if claims == nil {
		return "", errInvalidRefreshToken
	}

	// Check if the refresh token is revoked.
	isRevoked, err := postgres.IsRefreshTokenRevoked(refreshToken)
	if err != nil || isRevoked {
		return "", errInvalidRefreshToken
	}

	// Retrieve the user ID from the refresh token.
	userId, err := postgres.GetIdFromRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Generate a new access token for the user.
	return generateAccessToken(userId)
}

// logout invalidates the given refresh token.
//
// Returns an error if the refresh token is invalid, revoked, or if revocation fails.
func logout(refreshToken string) error {
	// Parse the claims from the refresh token.
	claims := parseTokenClaims(refreshToken)
	if claims == nil {
		return errInvalidRefreshToken
	}

	// Check if the refresh token is revoked.
	isRevoked, err := postgres.IsRefreshTokenRevoked(refreshToken)
	if err != nil || isRevoked {
		return errInvalidRefreshToken
	}

	// Revoke the refresh token in the database.
	return postgres.RevokeRefreshToken(refreshToken)
}
