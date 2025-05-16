package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationError represents a validation error
type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// Validate validates a struct and returns validation errors
func Validate(s interface{}) []ValidationError {
	var errors []ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = strings.ToLower(err.Field())
			element.Tag = err.Tag()
			element.Value = fmt.Sprintf("%v", err.Value())
			errors = append(errors, element)
		}
	}

	return errors
}

// FormatValidationErrors formats validation errors into a user-friendly message
func FormatValidationErrors(errors []ValidationError) string {
	if len(errors) == 0 {
		return ""
	}

	var messages []string
	for _, err := range errors {
		var msg string
		switch err.Tag {
		case "required":
			msg = fmt.Sprintf("%s is required", err.Field)
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters long", err.Field, err.Value)
		case "max":
			msg = fmt.Sprintf("%s must not exceed %s characters", err.Field, err.Value)
		case "gt":
			msg = fmt.Sprintf("%s must be greater than %s", err.Field, err.Value)
		case "gte":
			msg = fmt.Sprintf("%s must be greater than or equal to %s", err.Field, err.Value)
		case "ltefield":
			msg = fmt.Sprintf("%s must be before %s", err.Field, err.Value)
		default:
			msg = fmt.Sprintf("%s failed %s validation", err.Field, err.Tag)
		}
		messages = append(messages, msg)
	}

	return strings.Join(messages, "; ")
}
