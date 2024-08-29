package rest

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
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

func requirePermissions(resource, action string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(int)

		params := mux.Vars(r)

		var permissionCodes []string
		switch resource {
		case "collectionFilm":
			if action == "create" {
				permissionCodes = append(permissionCodes, "film"+":"+params["filmId"]+":"+"read")
				permissionCodes = append(permissionCodes, "collection"+":"+params["collectionId"]+":"+"update")
			} else if action == "read" {
				permissionCodes = append(permissionCodes, "collection"+":"+params["collectionId"]+":"+"read")
			} else {
				permissionCodes = append(permissionCodes, "collection"+":"+params["collectionId"]+":"+"update")
			}
		case "collection":
			if action == "create" {
				permissionCodes = append(permissionCodes, resource+":"+action)
			} else {
				permissionCodes = append(permissionCodes, resource+":"+params["collectionId"]+":"+action)
			}
		case "film":
			if action == "create" {
				permissionCodes = append(permissionCodes, resource+":"+action)
			} else {
				permissionCodes = append(permissionCodes, resource+":"+params["filmId"]+":"+action)
			}
		default:
			permissionCodes = append(permissionCodes, resource+":"+action)
		}

		permissions, err := postgres.GetUserPermissions(userId)
		if err != nil {
			handleDBError(w, r, err)
			return
		}

		for _, permissionCode := range permissionCodes {
			if !permissions.Include(permissionCode) {
				forbiddenResponse(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}
