package rest

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/metrics"
	"net/http"
	"strings"
	"time"
)

// Map of paths that do not require authentication.
var notRequireAuth = map[string]struct{}{
	"/":                     {},
	"/favicon.ico":          {},
	"/api":                  {},
	"/api/v1/healthcheck":   {},
	"/api/v1/auth/register": {},
	"/api/v1/auth/login":    {},
	"/api/v1/auth/refresh":  {},
}

var internalPaths = []string{
	"/swagger",
	"/metrics",
}

// authenticate ensures that requests have a valid authentication token or are to an excluded path.
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path

		// Check if the request path is in the list of paths that do not require authentication.
		if _, exists := notRequireAuth[requestPath]; exists {
			next.ServeHTTP(w, r)
			return
		}

		// Allow access to internal endpoints.
		if isInternalPath(requestPath) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract the token from the request header.
		tokenString := parseTokenFromHeader(r)
		if tokenString == "" {
			invalidAuthTokenResponse(w, r)
			return
		}

		// Parse the token to extract claims.
		claims := parseTokenClaims(tokenString)
		if claims == nil {
			invalidAuthTokenResponse(w, r)
			return
		}

		// Add the user ID from claims to the request context.
		ctx := context.WithValue(r.Context(), "userId", claims.UserId)
		r = r.WithContext(ctx)

		// Proceed to the next handler with the modified request.
		next.ServeHTTP(w, r)
	})
}

// requirePermissions ensures that the user has the necessary permissions for the specified resource and action.
func requirePermissions(resource, action string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(int)
		params := mux.Vars(r)

		// Initialize a slice to store the required permission codes.
		var permissionCodes []string

		// Determine permission codes based on the resource type and action.
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

		// Retrieve the user's permissions from the database.
		permissions, err := postgres.GetUserPermissions(userId)
		if err != nil {
			handleDBError(w, r, err)
			return
		}

		// Check if the user has all required permissions.
		for _, permissionCode := range permissionCodes {
			if !permissions.Include(permissionCode) {
				forbiddenResponse(w, r)
				return
			}
		}

		// Proceed to the next handler if permissions are valid.
		next.ServeHTTP(w, r)
	}
}

// logAndRecordMetrics logs the endpoint info and records metrics for the request.
func logAndRecordMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for internal endpoints
		if isInternalPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		sl.PrintEndpointInfo(r)
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()
		metrics.RecordRequestDuration(r, duration)
	})
}

func isInternalPath(requestPath string) bool {
	for _, path := range internalPaths {
		if strings.HasPrefix(requestPath, path) {
			return true
		}
	}
	return false
}
