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
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type envelope map[string]interface{}

func readJSON(p any, r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, data envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	e := json.NewEncoder(w)
	e.SetIndent("", "\t")

	err := e.Encode(data)
	if err != nil {
		sl.Log.Error(
			"failed to encode response data",
			slog.Any("error", err),
			slog.Any("request", r),
		)

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func parseIdParam(r *http.Request, paramName string) (int, error) {
	param := mux.Vars(r)[paramName]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func parseRequestBody(r *http.Request, v any) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errEmptyRequest
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	return nil
}

func parseTokenFromHeader(r *http.Request) string {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return ""
	}

	return strings.TrimPrefix(tokenHeader, "Bearer ")
}

func parseTokenClaims(tokenString string) *tokenClaims {
	claims := &tokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.TokenPassword), nil
	})

	if err != nil || !token.Valid {
		return nil
	}

	return claims
}

func addPermissionAndAssignToUser(userId, objectId int, objectType, action string) error {
	permission := fmt.Sprintf("%s:%d:%s", objectType, objectId, action)

	err := postgres.AddPermission(permission)
	if err != nil {
		return err
	}

	return postgres.AddUserPermissions(userId, permission)
}

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

func parseQueryString(qs url.Values, key string, defaultValue string) string {
	return parseQuery(qs, key, defaultValue, func(v string) (string, error) {
		return v, nil
	})
}

func parseQueryInt(qs url.Values, key string, defaultValue int) int {
	return parseQuery(qs, key, defaultValue, strconv.Atoi)
}

func parseQueryFloat(qs url.Values, key string, defaultValue float64) float64 {
	return parseQuery(qs, key, defaultValue, func(v string) (float64, error) {
		return strconv.ParseFloat(v, 64)
	})
}
