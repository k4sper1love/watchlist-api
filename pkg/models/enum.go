package models

// View status constants define the possible states of film viewing progress.
const (
	ViewStatusNotViewed  string = "not_viewed"  // The film has not been viewed.
	ViewStatusInProgress string = "in_progress" // The film is currently being watched.
	ViewStatusViewed     string = "viewed"      // The film has been fully watched.
)
