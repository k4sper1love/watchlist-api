package swagger

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/models"
)

type AuthResponse struct {
	User models.AuthResponse `json:"user"`
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
	CollectionFilms models.CollectionFilms `json:"collection_films"`
	Metadata        filters.Metadata       `json:"metadata"`
}
