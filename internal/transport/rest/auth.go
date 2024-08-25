package rest

import (
	"errors"
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

	accessToken, err := GenerateAccessToken(user.Id)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateAndSaveRefreshToken(user.Id)
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

	accessToken, err := GenerateAccessToken(user.Id)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateAndSaveRefreshToken(user.Id)
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
	isValid, err := postgres.IsRefreshTokenValid(refreshToken)
	if err != nil || !isValid {
		return "", errors.New("invalid or revoked refresh token")
	}

	userId, err := postgres.GetIdFromRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return GenerateAccessToken(userId)
}

func logout(refreshToken string) error {
	return postgres.RevokeRefreshToken(refreshToken)
}
