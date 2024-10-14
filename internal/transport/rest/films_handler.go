package rest

import (
	"database/sql"
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-api/pkg/validator"
	"net/http"
)

// AddFilm godoc
// @Summary Add new film
// @Description Add a new film. You will be granted the permissions to get, update, and delete it.
// @Tags films
// @Accept json
// @Produce json
// @Param film body swagger.FilmRequest true "Information about the new film".
// @Success 201 {object} swagger.FilmResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /films [post]
func addFilmHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var film models.Film

	if err := parseRequestBody(r, &film); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	film.UserID = userID

	if errs := validator.ValidateStruct(&film); errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	if err := postgres.AddFilm(&film); err != nil {
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

	writeJSON(w, r, http.StatusCreated, envelope{"film": film})
}

// GetFilm godoc
// @Summary Get film by ID
// @Description Get the film by ID. You must have permissions to get this film.
// @Tags films
// @Accept json
// @Produce json
// @Param film_id path int true "Film ID"
// @Success 200 {object} swagger.FilmResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /films/{film_id} [get]
func getFilmHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "filmID")
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

// GetFilms godoc
// @Summary Get user films
// @Description Get a list of films by user ID from authentication token. It also returns metadata.
// @Tags films
// @Accept json
// @Produce json
// @Param title query string false "Filter by `title`"
// @Param rating_min query number false "Filter by `minimum rating`"
// @Param rating_max query number false "Filter by `maximum rating`"
// @Param page query int false "Specify the desired `page`"
// @Param page_size query int false "Specify the desired `page size`"
// @Param sort query string false "Sorting by `id`, `title`, `rating`. Use `-` for desc"
// @Success 200 {object} swagger.FilmsResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /films [get]
func getFilmsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userID").(int)

	// Define an input structure to hold filter and pagination parameters.
	var input struct {
		Title     string
		MinRating float64
		MaxRating float64
		filters.Filters
	}

	// Parse query string parameters.
	qs := r.URL.Query()
	input.Title = parseQueryString(qs, "title", "")
	input.MinRating = parseQueryFloat(qs, "rating_min", 0)
	input.MaxRating = parseQueryFloat(qs, "rating_max", 0)
	input.Filters.Page = parseQueryInt(qs, "page", 1)
	input.Filters.PageSize = parseQueryInt(qs, "page_size", 5)
	input.Filters.Sort = parseQueryString(qs, "sort", "id")

	// Define safe sortable fields.
	input.Filters.SortSafeList = []string{
		"id", "title", "rating",
		"-id", "-title", "-rating",
	}

	if errs, err := filters.ValidateFilters(input.Filters); err != nil {
		serverErrorResponse(w, r, err)
		return
	} else if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	// Retrieve the list of films based on the filters.
	films, metadata, err := postgres.GetFilmsByUser(userId, input.Title, input.MinRating, input.MaxRating, input.Filters)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"films": films, "metadata": metadata})
}

// UpdateFilm godoc
// @Summary Update the film
// @Description Update the film by ID. You must have the permissions to update it.
// @Tags films
// @Accept json
// @Produce json
// @Param film_id path int true "Film ID"
// @Param film body swagger.FilmRequest true "New information about the film"
// @Success 200 {object} swagger.FilmResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 409 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /films/{film_id} [put]
func updateFilmHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "filmID")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	film, err := postgres.GetFilm(id)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	if err := parseRequestBody(r, film); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	film.ID = id

	if errs := validator.ValidateStruct(film); errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	if err := postgres.UpdateFilm(film); err != nil {
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

// DeleteFilm godoc
// @Summary Delete the film
// @Description Delete the film by ID. You must have the permissions to delete it.
// @Tags films
// @Accept json
// @Produce json
// @Param film_id path int true "Film ID"
// @Success 200 {object} swagger.MessageResponse
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 403 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /films/{film_id} [delete]
func deleteFilmHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "filmID")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// Verify that the film exists in the database.
	if _, err := postgres.GetFilm(id); err != nil {
		handleDBError(w, r, err)
		return
	}

	if err := postgres.DeleteFilm(id); err != nil {
		handleDBError(w, r, err)
		return
	}

	if err := deletePermissionCodes(id, "film"); err != nil {
		handleDBError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, envelope{"message": "film deleted"})
}
