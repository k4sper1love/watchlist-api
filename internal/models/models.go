/*
Package models defines the data structures used within the application.
These structures represent the main entities in the system, including users, collections,
films, and the relationships between them. The models include fields with JSON tags
for serialization and validation tags for input validation.

Each model contains fields relevant to its entity, including identifiers, attributes,
and timestamps for tracking creation and update times.
*/

package models

import (
	"time"
)

// User represents a user in the system.
type User struct {
	Id        int       `json:"id"`                                              // Unique identifier for the user.
	Username  string    `json:"username" validate:"required,username"`           // Username of the user; must be unique and valid.
	Email     string    `json:"email" validate:"required,email"`                 // Email address of the user; must be a valid email format.
	Password  string    `json:"password,omitempty" validate:"required,password"` // Password for the user account; omitted in responses for security.
	CreatedAt time.Time `json:"created_at"`                                      // Timestamp when the user was created.
	Version   int       `json:"-"`                                               // Internal version tracking; not included in JSON responses.
}

// Collection represents a collection of films created by a user.
type Collection struct {
	Id          int       `json:"id"`                                                 // Unique identifier for the collection.
	UserId      int       `json:"user_id"`                                            // Identifier of the user who created the collection.
	Name        string    `json:"name" validate:"required,min=3,max=100"`             // Name of the collection; required, between 3 and 100 characters.
	Description string    `json:"description,omitempty" validate:"omitempty,max=500"` // Description of the collection; optional, up to 500 characters.
	CreatedAt   time.Time `json:"created_at"`                                         // Timestamp when the collection was created.
	UpdatedAt   time.Time `json:"updated_at"`                                         // Timestamp when the collection was last updated.
}

// Film represents a film with its details and user-specific attributes.
type Film struct {
	Id          int       `json:"id"`                                                      // Unique identifier for the film.
	UserId      int       `json:"user_id"`                                                 // Identifier of the user who added the film.
	Title       string    `json:"title" validate:"required,min=3,max=100"`                 // Title of the film; required, between 3 and 100 characters.
	Year        int       `json:"year,omitempty" validate:"omitempty,gte=1888,lte=2100"`   // Release year of the film; optional, must be between 1888 and 2100.
	Genre       string    `json:"genre,omitempty" validate:"omitempty,alpha, max=100"`     // Genre of the film; optional, only alphabetic characters.
	Description string    `json:"description,omitempty" validate:"omitempty,max=1000"`     // Description of the film; optional, up to 1000 characters.
	Rating      float64   `json:"rating,omitempty" validate:"omitempty,gte=1,lte=10"`      // Rating of the film; optional, must be between 1 and 10.
	PhotoUrl    string    `json:"photo_url,omitempty" validate:"omitempty,url"`            // URL of the film's photo; optional, must be a valid URL.
	Comment     string    `json:"comment,omitempty" validate:"omitempty,max=500"`          // URL of the film's photo; optional, must be a valid URL.
	IsViewed    bool      `json:"is_viewed"`                                               // Indicates if the user has viewed the film.
	UserRating  float64   `json:"user_rating,omitempty" validate:"omitempty,gte=1,lte=10"` // User's rating of the film; optional, between 1 and 10.
	Review      string    `json:"review,omitempty" validate:"omitempty,max=500"`           // User's review of the film; optional, up to 500 characters.
	CreatedAt   time.Time `json:"created_at"`                                              // Timestamp when the film was added.
	UpdatedAt   time.Time `json:"updated_at"`                                              // Timestamp when the film details were last updated.
}

// CollectionFilm represents the association between a film and a collection.
type CollectionFilm struct {
	CollectionId int       `json:"collection_id"` // Identifier of the collection.
	FilmId       int       `json:"film_id"`       // Identifier of the film.
	CreatedAt    time.Time `json:"added_at"`      // Timestamp when the film was added to the collection.
	UpdatedAt    time.Time `json:"updated_at"`    // Timestamp when the association was last updated.
}
