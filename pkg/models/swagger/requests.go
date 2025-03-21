package swagger

import "time"

type LoginRequest struct {
	Username string `json:"username" example:"k4sper1love"`
	Password string `json:"password" example:"Secret1!"`
}

type RegisterRequest struct {
	LoginRequest
}

type UpdateUserRequest struct {
	Username string `json:"username" example:"new_username"`
	Email    string `json:"email" example:"new@example.com"`
}

type FilmRequest struct {
	IsFavorite  bool    `json:"is_favorite" example:"false"`
	Title       string  `json:"title" example:"My film"`
	Year        int     `json:"year" example:"2001"`
	Genre       string  `json:"genre" example:"Horror"`
	Description string  `json:"description" example:"This is description"`
	Rating      float64 `json:"rating" example:"6.7"`
	ImageURL    string  `json:"image_url" example:"http://k4sper1love.kz/images/default.png"`
	Comment     string  `json:"comment" example:"This is comment"`
	IsViewed    bool    `json:"is_viewed" example:"true"`
	UserRating  float64 `json:"user_rating" example:"5.5"`
	Review      string  `json:"review" example:"This is review."`
	URL         string  `json:"url" example:"https://www.kino.kz/film/689/"`
}

type CollectionRequest struct {
	IsFavorite  bool   `json:"is_favorite" example:"false"`
	Name        string `json:"name" example:"My collection"`
	Description string `json:"description,omitempty" example:"This is description"`
}

type CollectionFilmRequest struct {
	AddedAt time.Time `json:"added_at" example:"2024-09-04T13:37:24.87653+05:00"`
}
