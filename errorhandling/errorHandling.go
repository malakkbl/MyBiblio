package errorhandling

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	StatusCode int         `json:"-"`                 // HTTP status code
	Code       string      `json:"code"`              // Application-specific error code
	Message    string      `json:"message"`           // User-friendly error message
	Details    interface{} `json:"details,omitempty"` // Additional error details
	Debug      string      `json:"debug,omitempty"`   // Debug information (only in development)
}

// Standard error codes
const (
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeValidation         = "VALIDATION_ERROR"
	ErrCodeDatabase           = "DATABASE_ERROR"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeInvalidInput       = "INVALID_INPUT"
	ErrCodeInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrCodeDuplicateEntry     = "DUPLICATE_ENTRY"
	ErrCodeBadRequest         = "BAD_REQUEST"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeInvalidToken       = "INVALID_TOKEN"
	ErrCodeExpiredToken       = "EXPIRED_TOKEN"
	ErrCodeMissingToken       = "MISSING_TOKEN"
	ErrCodeWeakPassword       = "WEAK_PASSWORD"
	ErrCodeInvalidRole        = "INVALID_ROLE"
)

// Common application errors
var (
	ErrBookNotFound       = NewError(http.StatusNotFound, ErrCodeNotFound, "Book not found")
	ErrCustomerNotFound   = NewError(http.StatusNotFound, ErrCodeNotFound, "Customer not found")
	ErrAuthorNotFound     = NewError(http.StatusNotFound, ErrCodeNotFound, "Author not found")
	ErrOrderNotFound      = NewError(http.StatusNotFound, ErrCodeNotFound, "Order not found")
	ErrInsufficientStock  = NewError(http.StatusBadRequest, ErrCodeBadRequest, "Insufficient stock")
	ErrInvalidInput       = NewError(http.StatusBadRequest, ErrCodeBadRequest, "Invalid input")
	ErrInvalidCredentials = NewError(http.StatusUnauthorized, ErrCodeInvalidCredentials, "Invalid email or password")
	ErrInvalidToken       = NewError(http.StatusUnauthorized, ErrCodeInvalidToken, "Invalid authentication token")
	ErrExpiredToken       = NewError(http.StatusUnauthorized, ErrCodeExpiredToken, "Authentication token has expired")
	ErrMissingToken       = NewError(http.StatusUnauthorized, ErrCodeMissingToken, "Authentication token is missing")
	ErrWeakPassword       = NewError(http.StatusBadRequest, ErrCodeWeakPassword, "Password does not meet security requirements")
	ErrInvalidRole        = NewError(http.StatusBadRequest, ErrCodeInvalidRole, "Invalid user role specified")
)

// Error implements the error interface
func (e ErrorResponse) Error() string {
	return e.Message
}

// NewError creates a new ErrorResponse
func NewError(statusCode int, code string, message string) ErrorResponse {
	return ErrorResponse{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// WithDetails adds details to an ErrorResponse
func (e ErrorResponse) WithDetails(details interface{}) ErrorResponse {
	e.Details = details
	return e
}

// WithDebug adds debug information to an ErrorResponse
func (e ErrorResponse) WithDebug(debug string) ErrorResponse {
	e.Debug = debug
	return e
}

// HandleError writes an error response to http.ResponseWriter
func HandleError(w http.ResponseWriter, err error) {
	var errResp ErrorResponse

	switch e := err.(type) {
	case ErrorResponse:
		errResp = e
	case *json.SyntaxError:
		errResp = NewError(http.StatusBadRequest, ErrCodeInvalidInput, "Invalid JSON format")
	case *json.UnmarshalTypeError:
		errResp = NewError(http.StatusBadRequest, ErrCodeInvalidInput, "Invalid data type in JSON")
	default:
		errResp = NewError(http.StatusInternalServerError, ErrCodeInternalServer, "Internal server error")
	}

	// Set status code and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResp.StatusCode)
	json.NewEncoder(w).Encode(errResp)
}

// NewValidationError creates a validation error response
func NewValidationError(details interface{}) ErrorResponse {
	return NewError(http.StatusBadRequest, ErrCodeValidation, "Validation failed").
		WithDetails(details)
}

// NewDatabaseError creates a database error response
func NewDatabaseError(err error) ErrorResponse {
	return NewError(http.StatusInternalServerError, ErrCodeDatabase, "Database operation failed").
		WithDebug(err.Error())
}

// NewNotFoundError creates a not found error response
func NewNotFoundError(resource string, id interface{}) ErrorResponse {
	return NewError(
		http.StatusNotFound,
		ErrCodeNotFound,
		fmt.Sprintf("%s with ID %v not found", resource, id),
	)
}

// IsDuplicateKeyError checks if an error is a duplicate key error
func IsDuplicateKeyError(err error) bool {
	return errors.Is(err, ErrDuplicateKey)
}

var ErrDuplicateKey = errors.New("duplicate key error")
