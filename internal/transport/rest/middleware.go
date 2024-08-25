package rest

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"net/http"
)

func jwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{
			"/api/v1/healthcheck",
			"/api/v1/auth/register",
			"/api/v1/auth/login",
			"/api/v1/auth/refresh",
		}
		requestPath := r.URL.Path

		for _, path := range notAuth {
			if path == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenString := getTokenFromHeader(r)
		if tokenString == "" {
			invalidAuthTokenResponse(w, r)
			return
		}

		claims := &models.TokenClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.TokenPassword), nil
		})

		if err != nil || !token.Valid {
			invalidAuthTokenResponse(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
