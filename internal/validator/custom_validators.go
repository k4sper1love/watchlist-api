package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"unicode"
)

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
