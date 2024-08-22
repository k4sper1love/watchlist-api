package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"log"
	"net/http"
)

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addUserHandler serving:", r.URL.Path, r.Host)

	var user models.User
	err := parseRequestBody(r, &user)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	err = postgres.AddUser(&user)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"user": user})
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getUserHandler serving:", r.URL.Path, r.Host)
	id, err := parseIdParam(r, "userId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user, err := postgres.GetUserById(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"user": user})
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getUsersHandler serving:", r.URL.Path, r.Host)

	users, err := postgres.GetAllUsers()
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"users": users})
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updateUserHandler serving:", r.URL.Path, r.Host)
	id, err := parseIdParam(r, "userId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	_, err = postgres.GetUserById(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	var user models.User
	err = parseRequestBody(r, &user)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	user.Id = id

	err = postgres.UpdateUser(&user)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"user": user})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteUserHandler serving:", r.URL.Path, r.Host)
	id, err := parseIdParam(r, "userId")
	if err != nil {
		badRequestResponse(w, r, err)
	}

	_, err = postgres.GetUserById(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = postgres.DeleteUser(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "user deleted"})

}
