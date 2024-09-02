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

// init registers custom validation functions for "password" and "username" fields.
// This function runs automatically when the package is imported.
func init() {
	err := validate.RegisterValidation("password", passwordValidator)
	if err != nil {
		panic(err)
	}

	err = validate.RegisterValidation("username", usernameValidator)
	if err != nil {
		panic(err)
	}
}

// ValidateStruct validates the fields of a struct according to the registered rules.
// It returns a map of field names to their corresponding validation error messages if validation fails.
func ValidateStruct(s interface{}) map[string]string {
	// Validate the struct passed as an argument.
	err := validate.Struct(s)

	// If validation errors are encountered, they are handled below.
	if err != nil {
		var validationErrors validator.ValidationErrors

		// Check if the error is a validation error using the errors.As function.
		if errors.As(err, &validationErrors) {
			errs := make(map[string]string)

			// Iterate over each field error to extract and store the error message.
			for _, fieldErr := range validationErrors {
				// Convert the field name to snake_case.
				field := toSnakeCase(fieldErr.Field())
				// Retrieve and store the validation message.
				errs[field] = getValidationMessage(fieldErr)
			}

			// Return the map containing the validation errors.
			return errs
		}
	}

	// Return nil if no validation errors are found.
	return nil
}

// toSnakeCase converts a camelCase string to snake_case for easier readability.
func toSnakeCase(s string) string {
	// Compile a regular expression to match camelCase patterns in the string.
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	// Replace camelCase patterns with snake_case equivalents.
	snake := re.ReplaceAllString(s, "${1}_${2}")
	// Convert the resulting string to lowercase.
	return strings.ToLower(snake)
}
