package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"log"
	"net/http"
)

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getUserHandler serving:", r.URL.Path, r.Host)

	userId := r.Context().Value("userId").(int)

	user, err := postgres.GetUserById(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}
	user.Password = ""

	writeJSON(w, r, http.StatusOK, envelope{"user": user})
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updateUserHandler serving:", r.URL.Path, r.Host)

	userId := r.Context().Value("userId").(int)

	user, err := postgres.GetUserById(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = parseRequestBody(r, user)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	errs := models.ValidateStruct(user)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}
	user.Id = userId
	user.Password = ""

	err = postgres.UpdateUser(user)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"user": user})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteUserHandler serving:", r.URL.Path, r.Host)

	userId := r.Context().Value("userId").(int)

	_, err := postgres.GetUserById(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = postgres.DeleteUser(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "user deleted"})
}
