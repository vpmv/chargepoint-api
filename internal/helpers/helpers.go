package helpers

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func ValidateStruct[T any](v T) error {
	return validate.Struct(v)
}
