package models

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
	"strings"
	"unicode"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	err := validate.RegisterValidation("password", passwordValidator)
	if err != nil {
		log.Fatal(err)
	}

	err = validate.RegisterValidation("username", usernameValidator)
	if err != nil {
		log.Fatal(err)
	}
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

func toSnakeCase(s string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errs := make(map[string]string)
			for _, fieldErr := range validationErrors {
				field := toSnakeCase(fieldErr.Field())
				tag := fieldErr.Tag()

				switch tag {
				case "required":
					errs[field] = "is required field"
				case "username":
					errs[field] = "must contain only letters, numbers, dots, and underscores"
				case "email":
					errs[field] = "must be a valid email address"
				case "password":
					errs[field] = "must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special character"
				case "alphanum":
					errs[field] = "must contain only letters and numbers"
				case "alpha":
					errs[field] = "must contain only alphabetic characters"
				case "url":
					errs[field] = "must be a valid URL"
				case "min":
					errs[field] = "must be at least " + fieldErr.Param() + " characters long"
				case "max":
					errs[field] = "must be at most " + fieldErr.Param() + " character long"
				case "lte":
					errs[field] = "must be less than or equal to " + fieldErr.Param()
				case "gte":
					errs[field] = "must be greater than or equal to " + fieldErr.Param()
				case "jwt":
					errs[field] = "is not a valid JWT token"
				default:
					errs[field] = "invalid " + field + " value"
				}
			}
			return errs
		}
	}
	return nil
}
