package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/models"
	"log"
	"net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("registerHandler serving:", r.URL.Path, r.Host)

	var user models.User
	err := parseRequestBody(r, &user)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	resp, err := register(&user)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"user": resp})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("loginHandler serving:", r.URL.Path, r.Host)

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := parseRequestBody(r, &input)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	resp, err := login(input.Email, input.Password)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"user": resp})
}

func refreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := getTokenFromHeader(r)
	if refreshToken == "" {
		invalidAuthTokenResponse(w, r)
		return
	}

	newAccessToken, err := refreshAccessToken(refreshToken)
	if err != nil {
		invalidAuthTokenResponse(w, r)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"access_token": newAccessToken})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := getTokenFromHeader(r)
	if refreshToken == "" {
		invalidAuthTokenResponse(w, r)
		return
	}

	err := logout(refreshToken)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "token revoked"})
}