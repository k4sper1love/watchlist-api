package swagger

import "github.com/k4sper1love/watchlist-api/internal/models"

type ErrorResponse struct {
	Error string `json:"error" example:"some kind of error"`
}

type UserResponse struct {
	User models.User `json:"user"`
}
