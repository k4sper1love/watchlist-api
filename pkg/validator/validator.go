package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
	"strings"
)

func New() (*validator.Validate, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.RegisterValidation("password", passwordValidator)
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation("username", usernameValidator)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return validate, nil
}

func ValidateStruct(v *validator.Validate, s interface{}) map[string]string {
	err := v.Struct(s)

	if err != nil {
		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			errs := make(map[string]string)

			for _, fieldErr := range validationErrors {
				field := toSnakeCase(fieldErr.Field())
				errs[field] = getValidationMessage(fieldErr)
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
