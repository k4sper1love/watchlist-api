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

func addCollectionHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	userId := r.Context().Value("userId").(int)

	var collection models.Collection
	err := parseRequestBody(r, &collection)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	collection.UserId = userId

	v, err := validator.New()
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	errs := validator.ValidateStruct(v, &collection)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	err = postgres.AddCollection(&collection)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	actions := []string{"read", "update", "delete"}
	for _, action := range actions {
		err = addPermissionAndAssignToUser(userId, collection.Id, "collection", action)
		if err != nil {
			serverErrorResponse(w, r, err)
			return
		}
	}

	writeJSON(w, r, http.StatusCreated, envelope{"collection": collection})
}

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

func getCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	userId := r.Context().Value("userId").(int)

	var input struct {
		Name string
		filters.Filters
	}

	qs := r.URL.Query()

	input.Name = parseQueryString(qs, "name", "")

	input.Filters.Page = parseQueryInt(qs, "page", 1)
	input.Filters.PageSize = parseQueryInt(qs, "page_size", 5)

	input.Filters.Sort = parseQueryString(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		"id", "name", "created_at",
		"-id", "-name", "-created_at",
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

	collections, metadata, err := postgres.GetCollections(userId, input.Name, input.Filters)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collections": collections, "metadata": metadata})
}

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
	collection.Id = collectionId

	v, err := validator.New()
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	errs := validator.ValidateStruct(v, collection)
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

func deleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	collectionId, err := parseIdParam(r, "collectionId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

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

	writeJSON(w, r, http.StatusOK, envelope{"message": "collection deleted"})
}
