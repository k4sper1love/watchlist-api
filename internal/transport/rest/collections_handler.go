package rest

import (
	"database/sql"
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/validator"
	"net/http"
)

// AddCollection godoc
// @Summary Add new collection
// @Description Add a new collection. You will be granted the permissions to get, update, and delete it.
// @Tags collections
// @Accept json
// @Produce json
// @Param collection body swagger.CollectionRequest true "Information about the new collection"
// @Success 201 {object} swagger.CollectionResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections [post]
func addCollectionHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int)
	var collection models.Collection

	if err := parseRequestBody(r, &collection); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	collection.UserId = userId

	if errs := validator.ValidateStruct(&collection); errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	if err := postgres.AddCollection(&collection); err != nil {
		handleDBError(w, r, err)
		return
	}

	// Define permissions for the collection.
	actions := []string{"read", "update", "delete"}
	for _, action := range actions {
		if err := addPermissionAndAssignToUser(userId, collection.Id, "collection", action); err != nil {
			serverErrorResponse(w, r, err)
			return
		}
	}

	writeJSON(w, r, http.StatusCreated, envelope{"collection": collection})
}

// GetCollection godoc
// @Summary Get collection by ID
// @Description Get the collection by ID. You must have permissions to get this collection.
// @Tags collections
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Success 200 {object} swagger.CollectionResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id} [get]
func getCollectionHandler(w http.ResponseWriter, r *http.Request) {
	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	collection, err := postgres.GetCollection(collectionId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection": collection})
}

// GetCollections godoc
// @Summary Get user collections
// @Description Get a list of collections by user ID from authentication token. It also returns metadata.
// @Tags collections
// @Accept json
// @Produce json
// @Param name query string false "Filter by `name`"
// @Param page query int false "Specify the desired `page`"
// @Param page_size query int false "Specify the desired `page size`"
// @Param sort query string false "Sorting by `id`, `name`, `created_at`. Use `-` for desc"
// @Success 200 {object} swagger.CollectionsResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections [get]
func getCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int)

	// Define an input structure to hold filter and pagination parameters.
	var input struct {
		Name string
		filters.Filters
	}

	// Parse query string parameters.
	qs := r.URL.Query()
	input.Name = parseQueryString(qs, "name", "")
	input.Filters.Page = parseQueryInt(qs, "page", 1)
	input.Filters.PageSize = parseQueryInt(qs, "page_size", 5)
	input.Filters.Sort = parseQueryString(qs, "sort", "id")

	// Define safe sortable fields.
	input.Filters.SortSafeList = []string{
		"id", "name", "created_at",
		"-id", "-name", "-created_at",
	}

	if errs, err := filters.ValidateFilters(input.Filters); err != nil {
		serverErrorResponse(w, r, err)
		return
	} else if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	// Retrieve the list of collections based on the filters.
	collections, metadata, err := postgres.GetCollections(userId, input.Name, input.Filters)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collections": collections, "metadata": metadata})
}

// UpdateCollection godoc
// @Summary Update the collection
// @Description Update the collection by ID. You must have the permissions to update it.
// @Tags collections
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Param film body swagger.CollectionRequest true "New information about the collection"
// @Success 200 {object} swagger.FilmResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id} [put]
func updateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	collection, err := postgres.GetCollection(collectionId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	if err = parseRequestBody(r, collection); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	collection.Id = collectionId

	if errs := validator.ValidateStruct(collection); errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	if err := postgres.UpdateCollection(collection); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			editConflictResponse(w, r)
		default:
			handleDBError(w, r, err)
		}
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection": collection})
}

// DeleteCollection godoc
// @Summary Delete the collection
// @Description Delete the collection by ID. You must have the permissions to delete it.
// @Tags collections
// @Accept json
// @Produce json
// @Param collection_id path int true "Collection ID"
// @Success 200 {object} swagger.MessageResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /collections/{collection_id} [delete]
func deleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// Verify that the collection exists in the database.
	if _, err := postgres.GetCollection(collectionId); err != nil {
		handleDBError(w, r, err)
		return
	}

	if err := postgres.DeleteCollection(collectionId); err != nil {
		handleDBError(w, r, err)
		return
	}

	if err := deletePermissionCodes(collectionId, "collection"); err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "collection deleted"})
}
