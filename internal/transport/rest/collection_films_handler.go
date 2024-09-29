package rest

import (
	"database/sql"
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"net/http"
)

// AddCollectionFilm godoc
// @Summary Add film to collection
// @Description Add a film to the collection. You must have rights to get the film and update the collection.
// @Tags collectionFilms
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Param film_id path int true "Film ID"
// @Success 201 {object} swagger.CollectionResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id}/films/{film_id} [post]
func addCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
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
	collectionFilm := models.CollectionFilm{
		CollectionId: collectionId,
		FilmId:       filmId,
	}

	if err := postgres.AddCollectionFilm(&collectionFilm); err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"collection_film": collectionFilm})
}

// GetCollectionFilm godoc
// @Summary Get film from collection by ID
// @Description Get the film from collection by ID. You must have permissions to get this collection.
// @Tags collectionFilms
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Param film_id path int true "Film ID"
// @Success 200 {object} swagger.CollectionFilmResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id}/films/{film_id} [get]
func getCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
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

// GetCollectionFilms godoc
// @Summary Get films from collection
// @Description Get a list of films from collection by collection ID. It also returns metadata.
// @Description You must have permissions to get this collection.
// @Tags collectionFilms
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Param page query int false "Specify the desired `page`"
// @Param page_size query int false "Specify the desired `page size`"
// @Param sort query string false "Sorting by `film_id`, `added_at`. Use `-` for desc"
// @Success 200 {object} swagger.CollectionFilmsResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id}/films [get]
func getCollectionFilmsHandler(w http.ResponseWriter, r *http.Request) {
	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// Define an input structure to hold filter and pagination parameters.
	var input struct {
		filters.Filters
	}

	// Parse query string parameters.
	qs := r.URL.Query()
	input.Filters.Page = parseQueryInt(qs, "page", 1)
	input.Filters.PageSize = parseQueryInt(qs, "page_size", 5)
	input.Filters.Sort = parseQueryString(qs, "sort", "film_id")

	// Define safe sortable fields.
	input.Filters.SortSafeList = []string{
		"film_id", "added_at",
		"-film_id", "-added_at",
	}

	// Validate the filters.
	if errs, err := filters.ValidateFilters(input.Filters); err != nil {
		serverErrorResponse(w, r, err)
		return
	} else if errs != nil {
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

// UpdateCollectionFilm godoc
// @Summary Update film in collection
// @Description Update the film in the collection by ID`s. You must have the permissions to update collection.
// @Tags collectionFilms
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Param film_id path int true "Film ID"
// @Param film body swagger.CollectionFilmRequest true "New information about the film in the collection"
// @Success 200 {object} swagger.CollectionFilmResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id}/films/{film_id} [put]
func updateCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := parseRequestBody(r, collectionFilm); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	collectionFilm.CollectionId = collectionId
	collectionFilm.FilmId = filmId

	if err := postgres.UpdateCollectionFilm(collectionFilm); err != nil {
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

// DeleteCollectionFilms godoc
// @Summary Delete film from collection
// @Description Delete the film from the collection by ID. You must have the permissions to update collection.
// @Tags collectionFilms
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Param film_id path int true "Film ID"
// @Success 200 {object} swagger.MessageResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id}/films/{films_id} [delete]
func deleteCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
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
	if _, err := postgres.GetCollectionFilm(collectionId, filmId); err != nil {
		handleDBError(w, r, err)
		return
	}

	if err = postgres.DeleteCollectionFilm(collectionId, filmId); err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "collection_film deleted"})
}
