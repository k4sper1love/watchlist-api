package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/metrics"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// envelope is a map used for formatting JSON responses.
type envelope map[string]interface{}

// readJSON decodes JSON data from an io.Reader into the specified data structure.
func readJSON(target any, r io.Reader) error {
	return json.NewDecoder(r).Decode(target)
}

// writeJSON encodes the provided data structure into JSON and writes it to the http.ResponseWriter.
func writeJSON(w http.ResponseWriter, r *http.Request, status int, data envelope) {
	metrics.IncStatusCount(status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	e := json.NewEncoder(w)
	e.SetIndent("", "\t")

	if err := e.Encode(data); err != nil {
		sl.Log.Error("failed to encode response data", slog.Any("error", err), slog.Any("request", r))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// parseIdParam extracts an integer ID from a URL parameter.
func parseIdParam(r *http.Request, paramName string) (int, error) {
	param := mux.Vars(r)[paramName]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// parseRequestBody reads and decodes the JSON body of an HTTP request into the specified data structure.
func parseRequestBody(r *http.Request, target any) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errEmptyRequest
	}

	return json.Unmarshal(data, target)
}

// parseTokenFromHeader extracts the JWT token from the Authorization header of an HTTP request.
func parseTokenFromHeader(r *http.Request) string {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return ""
	}
	return strings.TrimPrefix(tokenHeader, "Bearer ")
}

// parseTokenClaims parses and validates a JWT token string, extracting the claims if valid.
func parseTokenClaims(tokenString string) *tokenClaims {
	claims := &tokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.TokenPass), nil
	})

	if err != nil || !token.Valid {
		return nil
	}
	return claims
}

// addPermissionAndAssignToUser adds a permission string to the database and assigns it to a user.
func addPermissionAndAssignToUser(userId, objectId int, objectType, action string) error {
	permission := fmt.Sprintf("%s:%d:%s", objectType, objectId, action)

	if err := postgres.AddPermission(permission); err != nil {
		return err
	}

	return postgres.AddUserPermissions(userId, permission)
}

// deletePermissionCodes removes permission codes related to an object from the database.
func deletePermissionCodes(objectId int, objectType string) error {
	codes := []string{
		fmt.Sprintf("%s:%d:%s", objectType, objectId, "read"),
		fmt.Sprintf("%s:%d:%s", objectType, objectId, "update"),
		fmt.Sprintf("%s:%d:%s", objectType, objectId, "delete"),
	}

	return postgres.DeletePermissions(codes...)
}

// parseQuery is a generic function for parsing query parameters from a URL.Values map.
func parseQuery[T any](qs url.Values, key string, defaultValue T, parseFunc func(string) (T, error)) T {
	value := qs.Get(key)

	if value == "" {
		return defaultValue
	}

	parsedValue, err := parseFunc(value)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}

// parseQueryString extracts a string query parameter from URL.Values.
func parseQueryString(qs url.Values, key string, defaultValue string) string {
	return parseQuery(qs, key, defaultValue, func(v string) (string, error) {
		return v, nil
	})
}

// parseQueryInt extracts an integer query parameter from URL.Values.
func parseQueryInt(qs url.Values, key string, defaultValue int) int {
	return parseQuery(qs, key, defaultValue, strconv.Atoi)
}

// parseQueryFloat extracts a float query parameter from URL.Values.
func parseQueryFloat(qs url.Values, key string, defaultValue float64) float64 {
	return parseQuery(qs, key, defaultValue, func(v string) (float64, error) {
		return strconv.ParseFloat(v, 64)
	})
}
