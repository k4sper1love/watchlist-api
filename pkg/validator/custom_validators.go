package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"unicode"
)

// validation Messages maps validation tags to human-readable error messages.
var validationMessages = map[string]string{
	"required": "is required field",
	"username": "must contain only letters, numbers, dots, and underscores",
	"email":    "must be a valid email address",
	"password": "must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special character",
	"alphanum": "must contain only letters and numbers",
	"alpha":    "must contain only alphabetic characters",
	"url":      "must be a valid URL",
	"lte":      "must be less than or equal to ",
	"gte":      "must be greater than or equal to ",
	"min":      "must be at least ",
	"max":      "must be at most ",
}

// getValidationMessage returns a human-readable error message for a given validation error.
func getValidationMessage(fe validator.FieldError) string {
	message, ok := validationMessages[fe.Tag()]
	if !ok {
		return "invalid field parameter"
	}

	switch fe.Tag() {
	case "lte", "gte":
		return message + fe.Param()
	case "min", "max":
		return message + fe.Param() + " characters long"
	default:
		return message
	}
}

// usernameValidator checks if a username contains only allowed characters: letters, numbers, dots, and underscores.
func usernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9._]+$`)
	return validUsername.MatchString(username)
}

// passwordValidator ensures a password meets security criteria:
// a minimum length of 8 characters, and includes upper, lower, number, and special characters.
func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	const minLength = 8
	hasMinLen := len(password) >= minLength
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	// Iterate through each character in the password to check for required types.
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		// Check if the character is a punctuation or symbol.
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Return true if all criteria are met.
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
