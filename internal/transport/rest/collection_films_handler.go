package rest

import (
	"database/sql"
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"net/http"
)

func addCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

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
	sl.PrintHandlerInfo(r)

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
	sl.PrintHandlerInfo(r)

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

	errs, err := filters.ValidateFilters(input.Filters)
	switch {
	case err != nil:
		serverErrorResponse(w, r, err)
		return
	case errs != nil:
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
	sl.PrintHandlerInfo(r)

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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			editConflictResponse(w, r)
		default:
			handleDBError(w, r, err)
		}
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection_film": collectionFilm})
}

func deleteCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

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
