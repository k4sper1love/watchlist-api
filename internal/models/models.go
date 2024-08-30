package models

import (
	"time"
)

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username" validate:"required,username"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password,omitempty" validate:"required,password"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

type Collection struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Name        string    `json:"name" validate:"required,min=3,max=100"`
	Description string    `json:"description,omitempty" validate:"omitempty,max=500"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Film struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Title       string    `json:"title" validate:"required,min=3,max=100"`
	Year        int       `json:"year,omitempty" validate:"omitempty,gte=1888,lte=2100"`
	Genre       string    `json:"genre,omitempty" validate:"omitempty,alpha"`
	Description string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	Rating      float64   `json:"rating,omitempty" validate:"omitempty,gte=1,lte=10"`
	PhotoUrl    string    `json:"photo_url,omitempty" validate:"omitempty,url"`
	Comment     string    `json:"comment,omitempty" validate:"omitempty,max=500"`
	IsViewed    bool      `json:"is_viewed"`
	UserRating  float64   `json:"user_rating,omitempty" validate:"omitempty,gte=1,lte=10"`
	Review      string    `json:"review,omitempty" validate:"omitempty,max=500"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CollectionFilm struct {
	CollectionId int       `json:"collection_id"`
	FilmId       int       `json:"film_id"`
	CreatedAt    time.Time `json:"added_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
