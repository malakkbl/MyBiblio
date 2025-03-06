package handlers

import (
	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/inmemorystores"
)

// SetupRouter initializes and returns the router (using `httprouter`)
func SetupRouter() *httprouter.Router {

	bookStore := inmemorystores.NewInMemoryBookStore()
	authorStore := inmemorystores.NewInMemoryAuthorStore()
	customerStore := inmemorystores.NewInMemoryCustomerStore()
	orderStore := inmemorystores.NewInMemoryOrderStore(bookStore)

	bookHandler := BookHandler{Store: bookStore}
	authorHandler := AuthorHandler{Store: authorStore}
	customerHandler := CustomerHandler{Store: customerStore}
	orderHandler := OrderHandler{Store: orderStore}

	router := httprouter.New()

	// Authentication Routes
	router.POST("/login", LoginUser)
	router.POST("/register", RegisterUser)

	// Books
	router.GET("/books/:id", bookHandler.GetBookByIDHandler)
	router.POST("/books", bookHandler.CreateBookHandler)
	router.PUT("/books/:id", bookHandler.UpdateBookHandler)
	router.DELETE("/books/:id", bookHandler.DeleteBookHandler)
	router.GET("/books", bookHandler.SearchBooksHandler)

	// Authors
	router.GET("/authors/:id", authorHandler.GetAuthorByIDHandler)
	router.POST("/authors", authorHandler.CreateAuthorHandler)
	router.PUT("/authors/:id", authorHandler.UpdateAuthorHandler)
	router.DELETE("/authors/:id", authorHandler.DeleteAuthorHandler)
	router.GET("/authors", authorHandler.ListAuthorsHandler)

	// Customers
	router.GET("/customers/:id", customerHandler.GetCustomerByIDHandler)
	router.POST("/customers", customerHandler.CreateCustomerHandler)
	router.PUT("/customers/:id", customerHandler.UpdateCustomerHandler)
	router.DELETE("/customers/:id", customerHandler.DeleteCustomerHandler)
	router.GET("/customers", customerHandler.ListCustomersHandler)

	// Orders
	router.GET("/orders", orderHandler.GetAllOrdersHandler)
	router.GET("/orders/:id", orderHandler.GetOrderByIDHandler)
	router.POST("/orders", orderHandler.CreateOrderHandler)
	router.PUT("/orders/:id", orderHandler.UpdateOrderHandler)
	router.DELETE("/orders/:id", orderHandler.DeleteOrderHandler)

	// Sales Reports
	router.GET("/sales-reports", GetSalesReportHandler)

	return router
}
