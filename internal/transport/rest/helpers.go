package rest

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
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
