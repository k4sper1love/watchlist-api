package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"log"
	"net/http"
)

func addUserFilm(w http.ResponseWriter, r *http.Request) {
	log.Println("addUserFilm serving:", r.URL.Path, r.Host)

	userId, err := parseIdParam(r, "userId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	var userFilm models.UserFilm
	err = parseRequestBody(r, &userFilm)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	userFilm.UserId = userId

	err = postgres.AddUserFilm(&userFilm)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"user_film": userFilm})
}

func getUserFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getUserFilmHandler serving:", r.URL.Path, r.Host)

	userId, err := parseIdParam(r, "userId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	filmId, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	userFilm, err := postgres.GetUserFilm(userId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"user_film": userFilm})
}

func getUserFilmsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getUserFilmsHandler serving:", r.URL.Path, r.Host)

	userId, err := parseIdParam(r, "userId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	userFilms, err := postgres.GetUserFilms(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"user_films": userFilms})
}

func updateUserFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updateUserFilm serving:", r.URL.Path, r.Host)

	userId, err := parseIdParam(r, "userId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	filmId, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	_, err = postgres.GetUserFilm(userId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	var userFilm models.UserFilm
	err = parseRequestBody(r, &userFilm)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	userFilm.UserId = userId
	userFilm.FilmId = filmId

	err = postgres.UpdateUserFilm(&userFilm)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"user_film": userFilm})
}

func deleteUserFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteUserFilm serving:", r.URL.Path, r.Host)

	userId, err := parseIdParam(r, "userId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	filmId, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	_, err = postgres.GetUserFilm(userId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = postgres.DeleteUserFilm(userId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "user_film deleted"})
}
