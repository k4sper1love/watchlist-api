package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type authResponse struct {
	*models.User
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func register(user *models.User) (*authResponse, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	err := postgres.AddUser(user)
	if err != nil {
		return nil, err
	}
	user.Password = ""

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

func login(email, password string) (*authResponse, error) {
	user, err := postgres.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	user.Password = ""

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

func refreshAccessToken(refreshToken string) (string, error) {
	claims := parseTokenClaims(refreshToken)
	if claims == nil {
		return "", errInvalidRefreshToken
	}

	isRevoked, err := postgres.IsRefreshTokenRevoked(refreshToken)
	if err != nil || isRevoked {
		return "", errInvalidRefreshToken
	}

	userId, err := postgres.GetIdFromRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return generateAccessToken(userId)
}

func logout(refreshToken string) error {
	claims := parseTokenClaims(refreshToken)
	if claims == nil {
		return errInvalidRefreshToken
	}

	isRevoked, err := postgres.IsRefreshTokenRevoked(refreshToken)
	if err != nil || isRevoked {
		return errInvalidRefreshToken
	}

	return postgres.RevokeRefreshToken(refreshToken)
}
