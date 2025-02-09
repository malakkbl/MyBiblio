package errorhandling

import "errors"

type ErrorResponse struct {
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

var ErrBookNotFound = errors.New("book not found")
var ErrCustomerNotFound = errors.New("customer not found")
var ErrAuthorNotFound = errors.New("author not found")
var ErrOrderNotFound = errors.New("order not found")
var ErrInsufficientStock = errors.New("insufficient stock")
