package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
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

	collectionFilms, err := postgres.GetCollectionFilms(collectionId)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"collection_films": collectionFilms})
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
