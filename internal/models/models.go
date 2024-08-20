package models

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type Film struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Year        int     `json:"year"`
	Genre       string  `json:"genre"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	PhotoUrl    string  `json:"photo_url"`
}

type Watchlist struct {
	UserId   int     `json:"user_id"`
	FilmId   int     `json:"film_id"`
	AddedAt  string  `json:"added_at"`
	IsViewed bool    `json:"is_viewed"`
	Rating   float64 `json:"rating"`
}
