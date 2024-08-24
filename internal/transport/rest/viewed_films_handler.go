package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"log"
	"net/http"
)

func addViewedFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addViewedFilmHandler serving:", r.URL.Path, r.Host)

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

	var viewedFilm models.ViewedFilm
	err = parseRequestBody(r, &viewedFilm)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	viewedFilm.UserId = userId
	viewedFilm.FilmId = filmId

	err = postgres.AddViewedFilm(&viewedFilm)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"viewed_film": viewedFilm})
}

func getViewedFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getViewedFilmHandler serving:", r.URL.Path, r.Host)

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

	viewedFilm, err := postgres.GetViewedFilm(userId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"viewed_film": viewedFilm})
}

func getViewedFilmsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getViewedFilmsHandler serving:", r.URL.Path, r.Host)

	userId, err := parseIdParam(r, "userId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	viewedFilms, err := postgres.GetViewedFilms(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"viewed_films": viewedFilms})
}

func updateViewedFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updateViewedFilmHandler serving:", r.URL.Path, r.Host)

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

	_, err = postgres.GetViewedFilm(userId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	var viewedFilm models.ViewedFilm
	err = parseRequestBody(r, &viewedFilm)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	viewedFilm.UserId = userId
	viewedFilm.FilmId = filmId

	err = postgres.UpdateViewedFilm(&viewedFilm)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"viewed_film": viewedFilm})
}

func deleteViewedFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteViewedFilmHandler serving:", r.URL.Path, r.Host)

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

	_, err = postgres.GetViewedFilm(userId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = postgres.DeleteViewedFilm(userId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "viewed_film deleted"})
}
