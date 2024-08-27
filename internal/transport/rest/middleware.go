package rest

import (
	"context"
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

		tokenString := parseTokenFromHeader(r)
		if tokenString == "" {
			invalidAuthTokenResponse(w, r)
			return
		}

		claims := parseTokenClaims(tokenString)
		if claims == nil {
			invalidAuthTokenResponse(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
