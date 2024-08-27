package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"log"
	"net/http"
)

func addFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addFilmHandler serving:", r.URL.Path, r.Host)

	userId := r.Context().Value("userId").(int)

	var film models.Film
	err := parseRequestBody(r, &film)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	film.UserId = userId

	errs := models.ValidateStruct(&film)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	err = postgres.AddFilm(&film)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"film": film})
}

func getFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getFilmHandler serving:", r.URL.Path, r.Host)

	id, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	film, err := postgres.GetFilm(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"film": film})
}

func getFilmsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getFilmsHandler serving:", r.URL.Path, r.Host)

	userId := r.Context().Value("userId").(int)

	films, err := postgres.GetFilmsByUser(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"films": films})
}

func updateFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updateFilmHandler serving:", r.URL.Path, r.Host)

	id, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	film, err := postgres.GetFilm(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = parseRequestBody(r, film)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	film.Id = id

	errs := models.ValidateStruct(film)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	err = postgres.UpdateFilm(film)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"film": film})
}

func deleteFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteFilmHandler serving:", r.URL.Path, r.Host)

	id, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	_, err = postgres.GetFilm(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = postgres.DeleteFilm(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "film deleted"})
}
