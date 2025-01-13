package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/metrics"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
		metrics.IncStatusCount(http.StatusInternalServerError)
	}
}

// parseIDParam extracts an integer ID from a URL parameter.
func parseIDParam(r *http.Request, paramName string) (int, error) {
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
func parseTokenClaims(tokenString, secret string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errInvalidToken
	}

	return claims, nil
}

// addPermissionAndAssignToUser adds a permission string to the database and assigns it to a user.
func addPermissionAndAssignToUser(userID, objectID int, objectType, action string) error {
	permission := fmt.Sprintf("%s:%d:%s", objectType, objectID, action)

	if err := postgres.AddPermission(permission); err != nil {
		return err
	}

	return postgres.AddUserPermissions(userID, permission)
}

// deletePermissionCodes removes permission codes related to an object from the database.
func deletePermissionCodes(objectID int, objectType string) error {
	codes := []string{
		fmt.Sprintf("%s:%d:%s", objectType, objectID, "read"),
		fmt.Sprintf("%s:%d:%s", objectType, objectID, "update"),
		fmt.Sprintf("%s:%d:%s", objectType, objectID, "delete"),
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

func parseQueryBool(qs url.Values, key string, defaultValue bool) bool {
	return parseQuery(qs, key, defaultValue, strconv.ParseBool)
}

func parseQueryBoolPtr(qs url.Values, key string) *bool {
	value := qs.Get(key)

	if value == "" {
		return nil
	}

	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}

	return &parsedValue
}

// generateString creates a random string of specified length from a set of characters.
func generateString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	str := make([]byte, length)
	for i := range str {
		str[i] = charset[r.Intn(len(charset))]
	}

	return string(str)
}

// generateUniqueUsername generates a unique username by combining a random string and a Telegram ID.
func generateUniqueUsername(length, telegramID int) string {
	for tries := 0; tries < 5; tries++ {
		base := generateString(length)
		username := fmt.Sprintf("%s_%d", base, telegramID)

		// Check if the username already exists in the database.
		if !postgres.IsUsernameExists(username) {
			return username
		}
	}

	// Return an empty string if no unique username is found after several attempts.
	return ""
}
