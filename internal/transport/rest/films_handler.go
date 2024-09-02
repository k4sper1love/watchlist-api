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

// addFilmHandler adds a new film to the database and assigns permissions to the user.
//
// Returns a JSON response with the created film or an error if the creation fails.
func addFilmHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Retrieve the user ID from the request context.
	userId := r.Context().Value("userId").(int)

	var film models.Film
	err := parseRequestBody(r, &film)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	film.UserId = userId // Assign the user ID to the film.

	errs := validator.ValidateStruct(&film)
	if errs != nil {
		failedValidationResponse(w, r, errs)
		return
	}

	err = postgres.AddFilm(&film)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	// Define permissions for the film.
	actions := []string{"read", "update", "delete"}
	for _, action := range actions {
		// Add permissions and assign them to the user.
		err = addPermissionAndAssignToUser(userId, film.Id, "film", action)
		if err != nil {
			serverErrorResponse(w, r, err)
			return
		}
	}

	writeJSON(w, r, http.StatusCreated, envelope{"film": film})
}

// getFilmHandler retrieves a film's details by its ID.
//
// Returns a JSON response with the film details or an error if retrieval fails.
func getFilmHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Parse the film ID from the request URL parameters.
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

// getFilmsHandler retrieves a list of films for the authenticated user with optional filters.
//
// Returns a JSON response with the list of films and metadata or an error if retrieval fails.
func getFilmsHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Retrieve the user ID from the request context.
	userId := r.Context().Value("userId").(int)

	// Define an input structure to hold filter and pagination parameters.
	var input struct {
		Title     string
		MinRating float64
		MaxRating float64
		filters.Filters
	}

	// Parse query string parameters from the request URL.
	qs := r.URL.Query()

	input.Title = parseQueryString(qs, "title", "")
	input.MinRating = parseQueryFloat(qs, "rating_min", 0)
	input.MaxRating = parseQueryFloat(qs, "rating_max", 0)

	input.Filters.Page = parseQueryInt(qs, "page", 1)
	input.Filters.PageSize = parseQueryInt(qs, "page_size", 5)

	input.Filters.Sort = parseQueryString(qs, "sort", "id")

	// Define a safe list of sortable fields.
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

	// Retrieve the list of films based on the filters.
	films, metadata, err := postgres.GetFilmsByUser(userId, input.Title, input.MinRating, input.MaxRating, input.Filters)
	if err != nil {
		handleDBError(w, r, err)
		return
	}

	// Respond with the list of films and metadata.
	writeJSON(w, r, http.StatusOK, envelope{"films": films, "metadata": metadata})
}

// updateFilmHandler updates the details of an existing film.
//
// Returns a JSON response with the updated film details or an error if the update fails.
func updateFilmHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	// Parse the film ID from the request URL parameters.
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
	film.Id = id // Ensure the film ID is set.

	errs := validator.ValidateStruct(film)
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

// deleteFilmHandler deletes a film from the database.
//
// Returns a JSON response confirming deletion or an error if the deletion fails.
func deleteFilmHandler(w http.ResponseWriter, r *http.Request) {
	sl.PrintHandlerInfo(r)

	id, err := parseIdParam(r, "filmId")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// Verify that the film exists in the database.
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

	// Confirm successful deletion with a JSON response.
	writeJSON(w, r, http.StatusOK, envelope{"message": "film deleted"})
}
