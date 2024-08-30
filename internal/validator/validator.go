package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
	"strings"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

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

func init() {
	err := validate.RegisterValidation("password", passwordValidator)
	if err != nil {
		log.Fatal(err)
	}

	err = validate.RegisterValidation("username", usernameValidator)
	if err != nil {
		log.Fatal(err)
	}
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
				message, exists := validationMessages[tag]

				if exists {
					if tag == "lte" || tag == "gte" {
						errs[field] = message + fieldErr.Param()
					} else if tag == "min" || tag == "max" {
						errs[field] = message + fieldErr.Param() + " characters long"
					} else {
						errs[field] = message
					}
				} else {
					errs[field] = "invalid " + field + " value"
				}
			}

			return errs
		}
	}

	return nil
}

func toSnakeCase(s string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}
