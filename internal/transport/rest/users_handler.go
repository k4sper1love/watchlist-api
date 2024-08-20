package rest

import (
	"encoding/json"
	"errors"
	"github.com/k4sper1love/wishlist-api/internal/database/postgres"
	"github.com/k4sper1love/wishlist-api/internal/models"
	"io"
	"log"
	"net/http"
)

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addUserHandler serving:", r.URL.Path, r.Host)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}
	if len(data) == 0 {
		BadRequestResponse(w, r, ErrEmptyRequest)
		return
	}

	var user models.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	err = postgres.AddUser(&user)
	if err != nil {
		if errors.Is(mapDBError(err), ErrAlreadyExists) {
			ConflictResponse(w, r)
			return
		}
		ServerErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusCreated, envelope{"user": user})
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getUserHandler serving:", r.URL.Path, r.Host)
	id, err := readIdParam(r)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	user, err := postgres.GetUserById(id)
	if err != nil {
		if errors.Is(mapDBError(err), ErrNotFound) {
			NotFoundResponse(w, r)
			return
		}
		ServerErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		ServerErrorResponse(w, r, err)
		return
	}
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getUsersHandler serving:", r.URL.Path, r.Host)

	users, err := postgres.GetAllUsers()
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"users": users})
	if err != nil {
		ServerErrorResponse(w, r, err)
		return
	}
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updateUserHandler serving:", r.URL.Path, r.Host)
	id, err := readIdParam(r)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	_, err = postgres.GetUserById(id)
	if err != nil {
		if errors.Is(mapDBError(err), ErrNotFound) {
			NotFoundResponse(w, r)
			return
		}
		ServerErrorResponse(w, r, err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}
	if len(data) == 0 {
		BadRequestResponse(w, r, ErrEmptyRequest)
		return
	}

	var user models.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}
	user.Id = id

	err = postgres.UpdateUser(&user)
	if err != nil {
		if errors.Is(mapDBError(err), ErrAlreadyExists) {
			ConflictResponse(w, r)
			return
		}
		ServerErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"user": user})
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteUserHandler serving:", r.URL.Path, r.Host)
	id, err := readIdParam(r)
	if err != nil {
		BadRequestResponse(w, r, err)
	}

	_, err = postgres.GetUserById(id)
	if err != nil {
		if errors.Is(mapDBError(err), ErrNotFound) {
			NotFoundResponse(w, r)
			return
		}
		ServerErrorResponse(w, r, err)
		return
	}

	err = postgres.DeleteUser(id)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"message": "user deleted"})
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}
