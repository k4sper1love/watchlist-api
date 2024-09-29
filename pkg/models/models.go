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

type User struct {
	Id        int       `json:"id" example:"1"`                                                       // Unique identifier for the user.
	Username  string    `json:"username" validate:"required,username" example:"john_doe"`             // Username of the user; must be unique and valid.
	Email     string    `json:"email" validate:"required,email" example:"john_doe@example.com"`       // Email address of the user; must be a valid email format.
	Password  string    `json:"password,omitempty" validate:"required,password" swaggerignore:"true"` // Password for the user account; omitted in responses for security.
	CreatedAt time.Time `json:"created_at" example:"2024-09-04T13:37:24.87653+05:00"`                 // Timestamp when the user was created.
	Version   int       `json:"-"`                                                                    // Internal version tracking; not included in JSON responses.
}

// Collection represents a collection of films created by a user.
type Collection struct {
	Id          int       `json:"id" example:"1"`                                                                   // Unique identifier for the collection.
	UserId      int       `json:"user_id" example:"1"`                                                              // Identifier of the user who created the collection.
	Name        string    `json:"name" validate:"required,min=3,max=100" example:"My collection"`                   // Name of the collection; required, between 3 and 100 characters.
	Description string    `json:"description,omitempty" validate:"omitempty,max=500" example:"This is description"` // Description of the collection; optional, up to 500 characters.
	CreatedAt   time.Time `json:"created_at" example:"2024-09-04T13:37:24.87653+05:00"`                             // Timestamp when the collection was created.
	UpdatedAt   time.Time `json:"updated_at" example:"2024-09-04T13:37:24.87653+05:00"`                             // Timestamp when the collection was last updated.
}

// Film represents a film with its details and user-specific attributes.
type Film struct {
	Id          int       `json:"id"  example:"1"`                                                                     // Unique identifier for the film.
	UserId      int       `json:"user_id" example:"1"`                                                                 // Identifier of the user who added the film.
	Title       string    `json:"title" validate:"required,min=3,max=100" example:"My film"`                           // Title of the film; required, between 3 and 100 characters.
	Year        int       `json:"year,omitempty" validate:"omitempty,gte=1888,lte=2100" example:"2001"`                // Release year of the film; optional, must be between 1888 and 2100.
	Genre       string    `json:"genre,omitempty" validate:"omitempty,alpha,max=100" example:"Horror"`                 // Genre of the film; optional, only alphabetic characters.
	Description string    `json:"description,omitempty" validate:"omitempty,max=1000" example:"This is description"`   // Description of the film; optional, up to 1000 characters.
	Rating      float64   `json:"rating,omitempty" validate:"omitempty,gte=1,lte=10" example:"6.7"`                    // Rating of the film; optional, must be between 1 and 10.
	PhotoUrl    string    `json:"photo_url,omitempty" validate:"omitempty,url" example:"https://placeimg.com/640/480"` // URL of the film's photo; optional, must be a valid URL.
	Comment     string    `json:"comment,omitempty" validate:"omitempty,max=500" example:"This is comment"`            // URL of the film's photo; optional, must be a valid URL.
	IsViewed    bool      `json:"is_viewed" example:"false"`                                                           // Indicates if the user has viewed the film.
	UserRating  float64   `json:"user_rating,omitempty" validate:"omitempty,gte=1,lte=10" example:"5.5"`               // User's rating of the film; optional, between 1 and 10.
	Review      string    `json:"review,omitempty" validate:"omitempty,max=500" example:"This is review"`              // User's review of the film; optional, up to 500 characters.
	CreatedAt   time.Time `json:"created_at" example:"2024-09-04T13:37:24.87653+05:00"`                                // Timestamp when the film was added.
	UpdatedAt   time.Time `json:"updated_at" example:"2024-09-04T13:37:24.87653+05:00"`                                // Timestamp when the film details were last updated.
}

// CollectionFilm represents the association between a film and a collection.
type CollectionFilm struct {
	CollectionId int       `json:"collection_id" example:"1"`                            // Identifier of the collection.
	FilmId       int       `json:"film_id" example:"1"`                                  // Identifier of the film.
	AddedAt      time.Time `json:"added_at" example:"2024-09-04T13:37:24.87653+05:00"`   // Timestamp when the film was added to the collection.
	UpdatedAt    time.Time `json:"updated_at" example:"2024-09-04T13:37:24.87653+05:00"` // Timestamp when the association was last updated.
}
