/*
Package validator provides functionality to validate struct fields using custom and built-in validation rules.
It supports registration of custom validation functions for fields like passwords and usernames.
This package also handles the conversion of validation errors to human-readable messages.
*/

package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

// validate is an instance of the validator library with required struct validation enabled.
var validate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	registerCustomValidators()
}

// registerCustomValidations registers custom validation functions.
func registerCustomValidators() {
	if err := validate.RegisterValidation("password", passwordValidator); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("username", usernameValidator); err != nil {
		panic(err)
	}
}

// ValidateStruct validates the fields of a struct according to the registered rules.
func ValidateStruct(s interface{}) map[string]string {
	// Validate the struct passed as an argument.
	if err := validate.Struct(s); err != nil {
		return extractValidationErrors(err)
	}
	return nil
}

// extractValidationErrors converts validation errors into a map of human-readable messages.
func extractValidationErrors(err error) map[string]string {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		errs := make(map[string]string)
		for _, fieldErr := range validationErrors {
			field := toSnakeCase(fieldErr.Field())
			errs[field] = getValidationMessage(fieldErr)
		}
		return errs
	}
	return nil
}

// toSnakeCase converts a camelCase string to snake_case for easier readability.
func toSnakeCase(s string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}
