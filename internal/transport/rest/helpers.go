package rest

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type envelope map[string]interface{}

func readIdParam(r *http.Request) (int, error) {
	param := mux.Vars(r)["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func readJSON(p interface{}, r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func writeJSON(w http.ResponseWriter, status int, data envelope) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	return e.Encode(data)
}
