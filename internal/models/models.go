package models

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type Film struct {
	Id          int     `json:"id"`
	CreatorId   int     `json:"creator_id"`
	Title       string  `json:"title"`
	Year        int     `json:"year"`
	Genre       string  `json:"genre"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	PhotoUrl    string  `json:"photo_url"`
}

type Collection struct {
	Id          int    `json:"id"`
	CreatorId   int    `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type UserFilm struct {
	UserId  int    `json:"user_id"`
	FilmId  int    `json:"film_id"`
	Comment string `json:"comment"`
	AddedAt string `json:"added_at"`
}

type UserViewed struct {
	UserId   int     `json:"user_id"`
	FilmId   int     `json:"film_id"`
	Rating   float64 `json:"rating"`
	Review   string  `json:"review"`
	ViewedAt string  `json:"viewed_at"`
}

type CollectionFilm struct {
	CollectionId int    `json:"collection_id"`
	FilmId       int    `json:"film_id"`
	AddedAt      string `json:"added_at"`
}

type UserCollection struct {
	UserId       int    `json:"user_id"`
	CollectionId int    `json:"collection_id"`
	AddedAt      string `json:"added_at"`
}
