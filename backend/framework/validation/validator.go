package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var msgs []string
	for _, e := range ve {
		msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}
	return strings.Join(msgs, "; ")
}

func Validate(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors ValidationErrors
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   err.Field(),
			Message: messageFor(err),
		})
	}

	return errors
}

func messageFor(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("minimum %s characters required", err.Param())
	case "max":
		return fmt.Sprintf("maximum %s characters allowed", err.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", err.Param())
	default:
		return fmt.Sprintf("validation failed on '%s' rule", err.Tag())
	}
}
