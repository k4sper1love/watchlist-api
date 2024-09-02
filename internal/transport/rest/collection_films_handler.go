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

// addCollectionFilmHandler adds a film to a collection.
//
// Returns a JSON response with the created collection-film relationship or an error if the addition fails.
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

	// Create a new CollectionFilm object with the parsed IDs.
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

// getCollectionFilmHandler retrieves the details of a film in a specific collection.
//
// Returns a JSON response with the collection-film relationship or an error if retrieval fails.
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

// getCollectionFilmsHandler retrieves a list of films in a specific collection with optional filters.
//
// Returns a JSON response with the list of collection-films and metadata or an error if retrieval fails.
func getCollectionFilmsHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// Define an input structure to hold filter and pagination parameters.
	var input struct {
		filters.Filters
	}

	// Parse query string parameters from the request URL.
	qs := r.URL.Query()

	input.Filters.Page = parseQueryInt(qs, "page", 1)
	input.Filters.PageSize = parseQueryInt(qs, "page_size", 5)

	input.Filters.Sort = parseQueryString(qs, "sort", "film_id")

	// Define a safe list of sortable fields.
	input.Filters.SortSafeList = []string{
		"film_id", "added_at",
		"-film_id", "-added_at",
	}

	// Validate the filters.
	errs, err := filters.ValidateFilters(input.Filters)
	switch {
	case err != nil:
		serverErrorResponse(w, r, err)
		return
	case errs != nil:
		failedValidationResponse(w, r, errs)
		return
	}

	// Retrieve the list of collection-films based on the filters.
	collectionFilms, metadata, err := postgres.GetCollectionFilms(collectionId, input.Filters)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection_films": collectionFilms, "metadata": metadata})
}

// updateCollectionFilmHandler updates the details of a film in a collection.
//
// Returns a JSON response with the updated collection-film relationship or an error if the update fails.
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
	collectionFilm.CollectionId = collectionId // Ensure the collection ID is set.
	collectionFilm.FilmId = filmId             // Ensure the film ID is set.

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

// deleteCollectionFilmHandler removes a film from a collection.
//
// Returns a JSON response confirming deletion or an error if the deletion fails.
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

	// Verify that the collection-film relationship exists in the database.
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
