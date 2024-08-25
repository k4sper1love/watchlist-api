package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Collection struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Film struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Title       string    `json:"title"`
	Year        int       `json:"year"`
	Genre       string    `json:"genre"`
	Description string    `json:"description"`
	Rating      float64   `json:"rating"`
	PhotoUrl    string    `json:"photo_url"`
	Comment     string    `json:"comment,omitempty"`
	IsViewed    bool      `json:"is_viewed"`
	UserRating  float64   `json:"user_rating,omitempty"`
	Review      string    `json:"review,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type CollectionFilm struct {
	CollectionId int       `json:"collection_id"`
	FilmId       int       `json:"film_id"`
	AddedAt      time.Time `json:"added_at"`
}

type TokenClaims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}
