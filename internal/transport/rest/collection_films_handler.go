package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-api/pkg/validator"
	"net/http"
)

// AddCollectionFilm godoc
// @Summary Add existing film to collection
// @Description Add existing film to the collection. You must have rights to get the film and update the collection.
// @Tags collectionFilms
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Param film_id path int true "Film ID"
// @Success 201 {object} swagger.CollectionFilmResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id}/films/{film_id} [post]
func addCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
	collectionID, err := parseIDParam(r, "collectionID")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	filmID, err := parseIDParam(r, "filmID")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// Create a new CollectionFilm object with the parsed IDs.
	collectionFilm := models.CollectionFilm{
		Collection: models.Collection{ID: collectionID},
		Film:       models.Film{ID: filmID},
	}

	if err := postgres.AddCollectionFilm(&collectionFilm); err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"collection_film": collectionFilm})
}

// addNewCollectionFilmHandler godoc
// @Summary Add new film and associate with collection
// @Description Create a new film and add it to the specified collection. You must have rights to create a film and update the collection.
// @Tags collectionFilms
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Param film body swagger.FilmRequest true "Information about the new film".
// @Success 201 {object} swagger.CollectionFilmResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id}/films [post]
func addNewCollectionFilmHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	collectionID, err := parseIDParam(r, "collectionID")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	var film models.Film
	if err := parseRequestBody(r, &film); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	film.UserID = userID

	setDefaultImage(r, &film)

	if errs := validator.ValidateStruct(&film); errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	if err := postgres.AddFilm(&film); err != nil {
		handleDBError(w, r, err)
		return
	}

	// Create a new CollectionFilm object with the parsed IDs.
	collectionFilm := models.CollectionFilm{
		Collection: models.Collection{ID: collectionID},
		Film:       models.Film{ID: film.ID},
	}

	if err := postgres.AddCollectionFilm(&collectionFilm); err != nil {
		handleDBError(w, r, err)
		return
	}

	// Define permissions for the film.
	actions := []string{"read", "update", "delete"}
	for _, action := range actions {
		if err := addPermissionAndAssignToUser(userID, film.ID, "film", action); err != nil {
			serverErrorResponse(w, r, err)
			return
		}
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
	collectionID, err := parseIDParam(r, "collectionID")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	filmID, err := parseIDParam(r, "filmID")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	collectionFilm := models.CollectionFilm{
		Collection: models.Collection{ID: collectionID},
		Film:       models.Film{ID: filmID},
	}

	if err := postgres.GetCollectionFilm(&collectionFilm); err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection_film": collectionFilm})
}

// GetCollectionFilms godoc
// @Summary Get films from collection
// @Description Retrieves a list of films from a specified collection. This includes pagination and sorting metadata.
// @Description You must have permissions to access this collection.
// @Tags collectionFilms
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Param page query int false "Specify the desired `page`"
// @Param page_size query int false "Specify the desired `page size`"
// @Param sort query string false "Sorting by `id`, `title`, `rating`. Use `-` for descending order"
// @Success 200 {object} swagger.CollectionFilmsResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id}/films [get]
func getCollectionFilmsHandler(w http.ResponseWriter, r *http.Request) {
	collectionID, err := parseIDParam(r, "collectionID")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	input, errs, err := parseAndValidateFilmsFilters(r)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	} else if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	collectionFilms := models.CollectionFilms{
		Collection: models.Collection{ID: collectionID},
	}

	metadata, err := postgres.GetCollectionFilms(&collectionFilms, input.Title, input.MinRating, input.MaxRating, input.Filters)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection_films": collectionFilms, "metadata": metadata})
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
	collectionID, err := parseIDParam(r, "collectionID")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	filmID, err := parseIDParam(r, "filmID")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	collectionFilm := models.CollectionFilm{
		Collection: models.Collection{ID: collectionID},
		Film:       models.Film{ID: filmID},
	}

	// Verify that the collection-film relationship exists in the database.
	if err := postgres.GetCollectionFilm(&collectionFilm); err != nil {
		handleDBError(w, r, err)
		return
	}

	if err = postgres.DeleteCollectionFilm(&collectionFilm); err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "collection_film deleted"})
}
