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
