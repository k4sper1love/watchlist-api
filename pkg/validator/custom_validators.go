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

// getValidationMessage returns the error message for a given validation error,
// customizing messages for "lte", "gte", "min", and "max" tags.
func getValidationMessage(fe validator.FieldError) string {
	message, ok := validationMessages[fe.Tag()]
	if !ok {
		return "invalid field param"
	}

	if fe.Tag() == "lte" || fe.Tag() == "gte" {
		return message + fe.Param()
	} else if fe.Tag() == "min" || fe.Tag() == "max" {
		return message + fe.Param() + "characters long"
	}

	return message
}

// usernameValidator checks if a username contains only allowed characters:
// letters, numbers, dots, and underscores.
func usernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9._]+$`)

	return validUsername.MatchString(username)
}

// passwordValidator ensures a password meets security criteria:
// minimum length of 8 characters, and contains upper, lower, number, and special characters.
func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var hasMinLen = len(password) >= 8 // Check if the password meets the minimum length of 8 characters.
	var hasUpper bool                  // Flag to track if the password contains an uppercase letter.
	var hasLower bool                  // Flag to track if the password contains a lowercase letter.
	var hasNumber bool                 // Flag to track if the password contains a numeric digit.
	var hasSpecial bool                // Flag to track if the password contains a special character

	// Iterate through each character in the password to set the respective flags.
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

	// Return true if all criteria are met, otherwise return false.
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
