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
// @Param film body swagger.FilmRequest true "Information about the new film".// @Success 201 {object} swagger.FilmResponse
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

	setDefaultImage(r, &film)

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
// @Param rating query string false "Filter by `rating`, can be a specific value or a range like 'min-max'"
// @Param year query string false "Filter by `year`"
// @Param user_rating query string false "Filter by `user_rating`"
// @Param is_viewed query bool false "Filter by `is_viewed` (true/false)"
// @Param is_favorite query bool false "Filter by `is_favorite` (true/false)"
// @Param has_url query bool false "Filter by `url` (true/false)"
// @Param exclude_collection query int false "Filter by `exclude collection`"
// @Param page query int false "Specify the desired `page`"
// @Param page_size query int false "Specify the desired `page size`"
// @Param sort query string false "Sorting by `id`, `title`, `rating`, `year`, `user_rating`, `is_viewed`. Use `-` for desc"
// @Success 200 {object} swagger.FilmsResponse
// @Failure 401 {object} swagger.ErrorResponse
// @Failure 422 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Security JWTAuth
// @Router /films [get]
func getFilmsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userID").(int)

	input, errs, err := parseAndValidateFilmsFilters(r)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	// Retrieve the list of films based on the filters.
	films, metadata, err := postgres.GetFilms(userId, input)
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
// @Param film body swagger.FilmRequest true "New information about the film"// @Success 200 {object} swagger.FilmResponse
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

	setDefaultImage(r, film)

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

// parseAndValidateFilmsFilters parses the incoming HTTP request for film filter and pagination parameters.
func parseAndValidateFilmsFilters(r *http.Request) (*models.FilmsQueryInput, map[string]string, error) {
	// Define an input structure to hold filter and pagination parameters.
	input := models.FilmsQueryInput{}
	// Parse query string parameters.
	qs := r.URL.Query()
	input.Title = parseQueryString(qs, "title", "")
	input.ExcludeCollection = parseQueryInt(qs, "exclude_collection", -1)
	input.Filters.Page = parseQueryInt(qs, "page", 1)
	input.Filters.PageSize = parseQueryInt(qs, "page_size", 5)
	input.Filters.Sort = parseQueryString(qs, "sort", "id")

	input.Rating = parseQueryString(qs, "rating", "")
	input.Year = parseQueryString(qs, "year", "")
	input.UserRating = parseQueryString(qs, "user_rating", "")

	isViewed := parseQueryBoolPtr(qs, "is_viewed")
	input.IsViewed = isViewed

	isFavorite := parseQueryBoolPtr(qs, "is_favorite")
	input.IsFavorite = isFavorite

	hasURL := parseQueryBoolPtr(qs, "has_url")
	input.HasURL = hasURL

	// Define safe sortable fields.
	input.Filters.SortSafeList = []string{
		"id", "title", "rating", "year", "is_viewed", "user_rating", "created_at",
		"-id", "-title", "-rating", "-year", "-is_viewed", "-user_rating", "-created_at",
	}

	errs, err := filters.ValidateFilters(input.Filters)

	return &input, errs, err
}
