package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/filters"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"log"
	"net/http"
)

func addCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addCollectionFilmHandler serving:", r.URL.Path, r.Host)

	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	filmId, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	var collectionFilm models.CollectionFilm
	//err = parseRequestBody(r, &collectionFilm)
	//if err != nil {
	//	badRequestResponse(w, r, err)
	//	return
	//}
	collectionFilm.CollectionId = collectionId
	collectionFilm.FilmId = filmId

	err = postgres.AddCollectionFilm(&collectionFilm)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"collection_film": collectionFilm})
}

func getCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getCollectionFilmHandler serving:", r.URL.Path, r.Host)

	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	filmId, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	collectionFilm, err := postgres.GetCollectionFilm(collectionId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection_film": collectionFilm})
}

func getCollectionFilmsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getCollectionFilmsHandler serving:", r.URL.Path, r.Host)

	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	var input struct {
		filters.Filters
	}

	qs := r.URL.Query()

	input.Filters.Page = parseQueryInt(qs, "page", 1)
	input.Filters.PageSize = parseQueryInt(qs, "page_size", 5)

	input.Filters.Sort = parseQueryString(qs, "sort", "film_id")

	input.Filters.SortSafeList = []string{
		"film_id", "added_at",
		"-film_id", "-added_at",
	}

	errs := filters.ValidateFilters(input.Filters)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	collectionFilms, metadata, err := postgres.GetCollectionFilms(collectionId, input.Filters)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection_films": collectionFilms, "metadata": metadata})
}

func updateCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updateCollectionFilmHandler serving:", r.URL.Path, r.Host)

	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	filmId, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	collectionFilm, err := postgres.GetCollectionFilm(collectionId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = parseRequestBody(r, collectionFilm)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	collectionFilm.CollectionId = collectionId
	collectionFilm.FilmId = filmId

	err = postgres.UpdateCollectionFilm(collectionFilm)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection_film": collectionFilm})
}

func deleteCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteCollectionFilmHandler serving:", r.URL.Path, r.Host)

	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	filmId, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	_, err = postgres.GetCollectionFilm(collectionId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = postgres.DeleteCollectionFilm(collectionId, filmId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "collection_film deleted"})
}
