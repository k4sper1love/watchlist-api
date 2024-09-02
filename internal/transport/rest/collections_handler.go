package rest

import (
	"database/sql"
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/validator"
	"net/http"
)

// addCollectionHandler adds a new collection to the database and assigns permissions to the user.
//
// Returns a JSON response with the created collection or an error if the creation fails.
func addCollectionHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Retrieve the user ID from the request context.
	userId := r.Context().Value("userId").(int)

	var collection models.Collection
	err := parseRequestBody(r, &collection)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	collection.UserId = userId // Assign the user ID to the collection.

	errs := validator.ValidateStruct(&collection)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	err = postgres.AddCollection(&collection)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	// Define permissions for the collection.
	actions := []string{"read", "update", "delete"}
	for _, action := range actions {
		// Add permissions and assign them to the user.
		err = addPermissionAndAssignToUser(userId, collection.Id, "collection", action)
		if err != nil {
			serverErrorResponse(w, r, err)
			return
		}
	}

	writeJSON(w, r, http.StatusCreated, envelope{"collection": collection})
}

// getCollectionHandler retrieves a collection's details by its ID.
//
// Returns a JSON response with the collection details or an error if retrieval fails.
func getCollectionHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

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

// getCollectionsHandler retrieves a list of collections for the authenticated user with optional filters.
//
// Returns a JSON response with the list of collections and metadata or an error if retrieval fails.
func getCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Retrieve the user ID from the request context.
	userId := r.Context().Value("userId").(int)

	// Define an input structure to hold filter and pagination parameters.
	var input struct {
		Name string
		filters.Filters
	}

	// Parse query string parameters from the request URL.
	qs := r.URL.Query()

	input.Name = parseQueryString(qs, "name", "")

	input.Filters.Page = parseQueryInt(qs, "page", 1)
	input.Filters.PageSize = parseQueryInt(qs, "page_size", 5)

	input.Filters.Sort = parseQueryString(qs, "sort", "id")

	// Define a safe list of sortable fields.
	input.Filters.SortSafeList = []string{
		"id", "name", "created_at",
		"-id", "-name", "-created_at",
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

	// Retrieve the list of collections based on the filters.
	collections, metadata, err := postgres.GetCollections(userId, input.Name, input.Filters)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collections": collections, "metadata": metadata})
}

// updateCollectionHandler updates the details of an existing collection.
//
// Returns a JSON response with the updated collection details or an error if the update fails.
func updateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

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

	err = parseRequestBody(r, collection)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	collection.Id = collectionId // Ensure the collection ID is set.

	errs := validator.ValidateStruct(collection)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	err = postgres.UpdateCollection(collection)
	if err != nil {
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

// deleteCollectionHandler deletes a collection from the database.
//
// Returns a JSON response confirming deletion or an error if the deletion fails.
func deleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// Verify that the collection exists in the database.
	_, err = postgres.GetCollection(collectionId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = postgres.DeleteCollection(collectionId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	// Confirm successful deletion with a JSON response.
	writeJSON(w, r, http.StatusOK, envelope{"message": "collection deleted"})
}
