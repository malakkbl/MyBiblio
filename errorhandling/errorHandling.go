package errorhandling

import "errors"

type ErrorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

var (
	ErrBookNotFound      = errors.New("book not found")
	ErrCustomerNotFound  = errors.New("customer not found")
	ErrAuthorNotFound    = errors.New("author not found")
	ErrOrderNotFound     = errors.New("order not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrValidation        = errors.New("validation error")
	ErrInvalidInput      = errors.New("invalid input")
)

// NewValidationError creates a new ErrorResponse for validation errors
func NewValidationError(details interface{}) ErrorResponse {
	return ErrorResponse{
		Message: "Validation failed",
		Details: details,
	}
}
