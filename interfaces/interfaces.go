package interfaces

import (
	"context"
	"time"

	"um6p.ma/finalproject/models"
)

type BookStore interface {
	CreateBook(ctx context.Context, book models.Book) (models.Book, error)
	GetBook(ctx context.Context, id int) (models.Book, error)
	UpdateBook(ctx context.Context, id int, book models.Book) (models.Book, error)
	DeleteBook(ctx context.Context, id int) error
	SearchBooks(ctx context.Context, criteria models.SearchCriteria) ([]models.Book, error)
	GetAllBooks(ctx context.Context) ([]models.Book, error)
}

type CustomerStore interface {
	CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error)
	GetCustomer(ctx context.Context, id int) (models.Customer, error)
	UpdateCustomer(ctx context.Context, id int, customer models.Customer) (models.Customer, error)
	DeleteCustomer(ctx context.Context, id int) error
	GetAllCustomers(ctx context.Context) ([]models.Customer, error)
}

type AuthorStore interface {
	CreateAuthor(ctx context.Context, author models.Author) (models.Author, error)
	GetAuthor(ctx context.Context, id int) (models.Author, error)
	UpdateAuthor(ctx context.Context, id int, author models.Author) (models.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	GetAllAuthors(ctx context.Context) ([]models.Author, error)
}

type OrderStore interface {
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	GetOrder(ctx context.Context, id int) (models.Order, error)
	UpdateOrder(ctx context.Context, id int, order models.Order) (models.Order, error)
	DeleteOrder(ctx context.Context, id int) error
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	GetOrdersInTimeRange(ctx context.Context, start, end time.Time) ([]models.Order, error)
}
