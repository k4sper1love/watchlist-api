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
	"github.com/golang-jwt/jwt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"time"
)

// JWTClaims defines the structure of JWT claims for user authentication by credentials.
type JWTClaims struct {
	Sub string `json:"sub"`
	jwt.StandardClaims
}

// Credentials represents the information required for user registration and authentication.
type Credentials struct {
	TelegramID int    `json:"telegram_id,omitempty"`
	Username   string `json:"username" validate:"required,username" example:"john_doe"`                  // Username of the user; must be unique and valid.
	Email      string `json:"email,omitempty" validate:"omitempty,email" example:"john_doe@example.com"` // Email address of the user; must be a valid email format.
	Password   string `json:"password,omitempty" validate:"required,password" swaggerignore:"true"`      // Password for the user account; omitted in responses for security.
}

// User represents the user data stored in the system.
type User struct {
	ID         int       `json:"id" example:"1"` // Unique identifier for the user.
	TelegramID int       `json:"telegram_id,omitempty" example:"123456789"`
	Username   string    `json:"username,omitempty" validate:"omitempty,username" example:"john_doe"`       // Username of the user; must be unique and valid.
	Email      string    `json:"email,omitempty" validate:"omitempty,email" example:"john_doe@example.com"` // Email address of the user; must be a valid email format.
	Password   string    `json:"password,omitempty" swaggerignore:"true"`                                   // Password for the user account; omitted in responses for security.
	CreatedAt  time.Time `json:"created_at" example:"2024-09-04T13:37:24.87653+05:00"`                      // Timestamp when the user was created.
	Version    int       `json:"-"`                                                                         // Internal version tracking; not included in JSON responses.
}

// AuthResponse represents the response returned upon successful authentication.
type AuthResponse struct {
	*User
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs.eyJzdWIilIn0.iTNuOHMObmeRmKU"` // JWT Access Token used to access protected resources.
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOI6IkpXVCJ9.eyJzdk5EbifQ.4CfEaMw6Ur_fszI"` // JWT Refresh Token used to obtain a new Access Token when it expires.
}

// Collection represents a collection of films created by a user.
type Collection struct {
	ID          int       `json:"id" example:"1"`      // Unique identifier for the collection.
	UserID      int       `json:"user_id" example:"1"` // Identifier of the user who created the collection.
	IsFavorite  bool      `json:"is_favorite" example:"false"`
	Name        string    `json:"name" validate:"required,min=3,max=100" example:"My collection"`                   // Name of the collection; required, between 3 and 100 characters.
	Description string    `json:"description,omitempty" validate:"omitempty,max=500" example:"This is description"` // Description of the collection; optional, up to 500 characters.
	TotalFilms  int       `json:"total_films" example:"5"`                                                          // Total number of films in the collection.
	CreatedAt   time.Time `json:"created_at" example:"2024-09-04T13:37:24.87653+05:00"`                             // Timestamp when the collection was created.
	UpdatedAt   time.Time `json:"updated_at" example:"2024-09-04T13:37:24.87653+05:00"`                             // Timestamp when the collection was last updated.
}

// Film represents a film with its details and user-specific attributes.
type Film struct {
	ID          int       `json:"id"  example:"1"`     // Unique identifier for the film.
	UserID      int       `json:"user_id" example:"1"` // Identifier of the user who added the film.
	IsFavorite  bool      `json:"is_favorite" example:"false"`
	Title       string    `json:"title" validate:"required,min=3,max=100" example:"My film"`                           // Title of the film; required, between 3 and 100 characters.
	Year        int       `json:"year,omitempty" validate:"omitempty,gte=1888,lte=2100" example:"2001"`                // Release year of the film; optional, must be between 1888 and 2100.
	Genre       string    `json:"genre,omitempty" validate:"omitempty,max=100" example:"Horror"`                       // Genre of the film; optional.
	Description string    `json:"description,omitempty" validate:"omitempty,max=1000" example:"This is description"`   // Description of the film; optional, up to 1000 characters.
	Rating      float64   `json:"rating,omitempty" validate:"omitempty,gte=1,lte=10" example:"6.7"`                    // Rating of the film; optional, must be between 1 and 10.
	ImageURL    string    `json:"image_url,omitempty" validate:"omitempty,url" example:"https://placeimg.com/640/480"` // URL of the film's image; optional, must be a valid URL.
	Comment     string    `json:"comment,omitempty" validate:"omitempty,max=500" example:"This is comment"`            // User's comment of the film; optional, up to 500 characters.
	IsViewed    bool      `json:"is_viewed" example:"true"`                                                            // Indicates if the user has viewed the film.
	UserRating  float64   `json:"user_rating,omitempty" validate:"omitempty,gte=1,lte=10" example:"5.5"`               // User's rating of the film; optional, between 1 and 10.
	Review      string    `json:"review,omitempty" validate:"omitempty,max=500" example:"This is review"`              // User's review of the film; optional, up to 500 characters.
	URL         string    `json:"url,omitempty" validate:"omitempty,url" example:"https://www.imdb.com/video"`         // URL for additional film information (e.g., IMDb or trailer); optional, must be valid.
	CreatedAt   time.Time `json:"created_at" example:"2024-09-04T13:37:24.87653+05:00"`                                // Timestamp when the film was added.
	UpdatedAt   time.Time `json:"updated_at" example:"2024-09-04T13:37:24.87653+05:00"`                                // Timestamp when the film details were last updated.
}

// CollectionFilm represents the association between a film and a collection.
type CollectionFilm struct {
	Collection Collection `json:"collection"`                                           // Identifier of the collection.
	Film       Film       `json:"film"`                                                 // Identifier of the film.
	AddedAt    time.Time  `json:"added_at" example:"2024-09-04T13:37:24.87653+05:00"`   // Timestamp when the film was added to the collection.
	UpdatedAt  time.Time  `json:"updated_at" example:"2024-09-04T13:37:24.87653+05:00"` // Timestamp when the association was last updated.
}

// CollectionFilms represents the association between a film and a collection.
type CollectionFilms struct {
	Collection Collection `json:"collection"` // Identifier of the collection.
	Films      []Film     `json:"films"`      // Identifier of the film.
}

// FilmsQueryInput holds the parameters for querying films, including title, rating range, and filter options.
type FilmsQueryInput struct {
	filters.Filters
	Title             string
	ExcludeCollection int
	Rating            string
	Year              string
	IsViewed          *bool
	UserRating        string
	HasURL            *bool
	IsFavorite        *bool
}
