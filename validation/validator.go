package validation

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validation functions
	validate.RegisterValidation("future_date", validateFutureDate)
	validate.RegisterValidation("past_date", validatePastDate)
	validate.RegisterValidation("valid_isbn", validateISBN)
	validate.RegisterValidation("valid_status", validateOrderStatus)
	validate.RegisterValidation("passwd", validatePassword)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message,omitempty"`
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
			element.Message = formatErrorMessage(element.Field, element.Tag, element.Value)
			errors = append(errors, element)
		}
	}

	return errors
}

// formatErrorMessage creates a user-friendly error message based on the validation error
func formatErrorMessage(field string, tag string, value string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, value)
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, value)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, value)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, value)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "ltefield":
		return fmt.Sprintf("%s must be before %s", field, value)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, value)
	case "valid_isbn":
		return fmt.Sprintf("%s must be a valid ISBN-10 or ISBN-13", field)
	case "valid_status":
		return fmt.Sprintf("%s must be a valid order status", field)
	default:
		return fmt.Sprintf("%s failed %s validation", field, tag)
	}
}

// Custom validation functions

// validateFutureDate ensures a date is in the future
func validateFutureDate(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return date.After(time.Now())
}

// validatePastDate ensures a date is in the past
func validatePastDate(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return date.Before(time.Now())
}

// validateISBN validates ISBN-10 and ISBN-13 formats
func validateISBN(fl validator.FieldLevel) bool {
	isbn := fl.Field().String()

	// Remove hyphens and spaces
	isbn = strings.ReplaceAll(isbn, "-", "")
	isbn = strings.ReplaceAll(isbn, " ", "")

	// Check length (ISBN-10 or ISBN-13)
	if len(isbn) != 10 && len(isbn) != 13 {
		return false
	}

	// For simplicity, we'll just check the length
	// In a real application, you would implement the full ISBN checksum algorithm
	return true
}

// validateOrderStatus validates order status values
func validateOrderStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
	}
	return validStatuses[strings.ToLower(status)]
}

// validatePassword ensures password meets security requirements
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

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

	return len(password) >= 8 && hasUpper && hasLower && hasNumber && hasSpecial
}

// AddCustomError adds a custom validation error
func AddCustomError(errors []ValidationError, field, tag, value, message string) []ValidationError {
	return append(errors, ValidationError{
		Field:   field,
		Tag:     tag,
		Value:   value,
		Message: message,
	})
}
