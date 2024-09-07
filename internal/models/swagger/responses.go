package swagger

import (
	"github.com/k4sper1love/watchlist-api/internal/models"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
)

type AuthResponse struct {
	*models.User
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs.eyJzdWIilIn0.iTNuOHMObmeRmKU"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOI6IkpXVCJ9.eyJzdk5EbifQ.4CfEaMw6Ur_fszI"`
}

type systemInfo struct {
	Environment string `json:"environment" example:"prod"`
	Version     string `json:"version" example:"2024-09-05T02:21:09Z-83438704281d314aae2e1ae-dirty"`
}
type HealthCheckResponse struct {
	Status     string     `json:"status" example:"enabled"`
	SystemInfo systemInfo `json:"system_info"`
}

type MessageResponse struct {
	Message string `json:"message" example:"some kind of success message"`
}
type ErrorResponse struct {
	Error string `json:"error" example:"some kind of error"`
}

type UserResponse struct {
	User models.User `json:"user"`
}

type AccessTokenResponse struct {
	Token string `json:"access_token" example:"eyJhbGciOI6IkpXVCJ9.eyJzdk5EbifQ.4CfEaMw6Ur_fszI"`
}

type FilmResponse struct {
	Film models.Film `json:"film"`
}

type FilmsResponse struct {
	Films    []models.Film    `json:"films"`
	Metadata filters.Metadata `json:"metadata"`
}

type CollectionResponse struct {
	Collection models.Collection `json:"collection"`
}

type CollectionsResponse struct {
	Collections []models.Collection `json:"collections"`
	Metadata    filters.Metadata    `json:"metadata"`
}

type CollectionFilmResponse struct {
	CollectionFilm models.CollectionFilm `json:"collection_film"`
}

type CollectionFilmsResponse struct {
	CollectionFilms []models.CollectionFilm `json:"collection_films"`
	Metadata        filters.Metadata        `json:"metadata"`
}
