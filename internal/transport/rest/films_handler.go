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

func addFilmHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	userId := r.Context().Value("userId").(int)

	var film models.Film
	err := parseRequestBody(r, &film)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	film.UserId = userId

	v, err := validator.New()
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	errs := validator.ValidateStruct(v, &film)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	err = postgres.AddFilm(&film)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	actions := []string{"read", "update", "delete"}
	for _, action := range actions {
		err = addPermissionAndAssignToUser(userId, film.Id, "film", action)
		if err != nil {
			serverErrorResponse(w, r, err)
			return
		}
	}

	writeJSON(w, r, http.StatusCreated, envelope{"film": film})
}

func getFilmHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	id, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	film, err := postgres.GetFilm(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"film": film})
}

func getFilmsHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	userId := r.Context().Value("userId").(int)

	var input struct {
		Title     string
		MinRating float64
		MaxRating float64
		filters.Filters
	}

	qs := r.URL.Query()

	input.Title = parseQueryString(qs, "title", "")
	input.MinRating = parseQueryFloat(qs, "rating_min", 0)
	input.MaxRating = parseQueryFloat(qs, "rating_max", 0)

	input.Filters.Page = parseQueryInt(qs, "page", 1)
	input.Filters.PageSize = parseQueryInt(qs, "page_size", 5)

	input.Filters.Sort = parseQueryString(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		"id", "title", "rating",
		"-id", "-title", "-rating",
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

	films, metadata, err := postgres.GetFilmsByUser(userId, input.Title, input.MinRating, input.MaxRating, input.Filters)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"films": films, "metadata": metadata})
}

func updateFilmHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	id, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	film, err := postgres.GetFilm(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = parseRequestBody(r, film)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	film.Id = id

	v, err := validator.New()
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	errs := validator.ValidateStruct(v, film)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	err = postgres.UpdateFilm(film)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			editConflictResponse(w, r)
		default:
			handleDBError(w, r, err)
		}
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"film": film})
}

func deleteFilmHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	id, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	_, err = postgres.GetFilm(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	err = postgres.DeleteFilm(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "film deleted"})
}
