package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"io"
	"log"
	"net/http"
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
		log.Println(r, err)
		w.WriteHeader(500)
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
