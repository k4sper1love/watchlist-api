package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"log"
	"net/http"
)

func addCollectionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addCollectionHandler serving:", r.URL.Path, r.Host)

	userId := r.Context().Value("userId").(int)

	var collection models.Collection
	err := parseRequestBody(r, &collection)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	collection.UserId = userId

	errs := models.ValidateStruct(&collection)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	err = postgres.AddCollection(&collection)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, envelope{"collection": collection})
}

func getCollectionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getCollectionHandler serving:", r.URL.Path, r.Host)

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
	log.Println("getCollectionsHandler serving:", r.URL.Path, r.Host)

	userId := r.Context().Value("userId").(int)

	collections, err := postgres.GetCollections(userId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collections": collections})
}

func updateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updateCollectionHandler serving:", r.URL.Path, r.Host)

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

	errs := models.ValidateStruct(collection)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	err = postgres.UpdateCollection(collection)
	if err != nil {
		handleDBError(w, r, err)
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection": collection})
}

func deleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteCollectionHandler serving:", r.URL.Path, r.Host)

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
