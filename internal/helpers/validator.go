package helpers

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

// ValidateEntity validates an entity of type T by its validation tags
func ValidateEntity[T any](v T) validator.ValidationErrors {
	var validationErrors validator.ValidationErrors = nil

	err := validate.Struct(v)
	if err != nil {
		errors.As(err, &validationErrors)
	}

	return validationErrors
}
