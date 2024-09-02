package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"unicode"
)

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

func usernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9._]+$`)

	return validUsername.MatchString(username)
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var hasMinLen = len(password) >= 8
	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
