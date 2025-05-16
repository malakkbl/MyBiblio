package handlers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/constants"
	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/inmemorystores"
	httputil "um6p.ma/finalproject/internal/http"
)

// SetupRouter initializes and returns the router
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

	// Custom error handler for router
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusInternalServerError,
			errorhandling.ErrCodeInternalServer,
			"Internal server error",
		).WithDebug(fmt.Sprint(err)))
	}

	// Public routes (no authentication required)
	router.POST("/login", LoginUser)
	router.POST("/register", RegisterUser)

	// Books routes
	router.GET("/books/:id", httputil.Wrap(bookHandler.GetBookByIDHandler, "read:books"))
	router.GET("/books", httputil.Wrap(bookHandler.SearchBooksHandler, "read:books"))
	router.POST("/books", httputil.WrapWithRoles(bookHandler.CreateBookHandler, constants.RoleAdmin, constants.RoleManager))
	router.PUT("/books/:id", httputil.WrapWithRoles(bookHandler.UpdateBookHandler, constants.RoleAdmin, constants.RoleManager))
	router.DELETE("/books/:id", httputil.WrapWithRole(bookHandler.DeleteBookHandler, constants.RoleAdmin))

	// Authors routes
	router.GET("/authors/:id", httputil.Wrap(authorHandler.GetAuthorByIDHandler, "read:authors"))
	router.GET("/authors", httputil.Wrap(authorHandler.ListAuthorsHandler, "read:authors"))
	router.POST("/authors", httputil.WrapWithRoles(authorHandler.CreateAuthorHandler, constants.RoleAdmin, constants.RoleManager))
	router.PUT("/authors/:id", httputil.WrapWithRoles(authorHandler.UpdateAuthorHandler, constants.RoleAdmin, constants.RoleManager))
	router.DELETE("/authors/:id", httputil.WrapWithRole(authorHandler.DeleteAuthorHandler, constants.RoleAdmin))

	// Customers routes
	router.GET("/customers/:id", httputil.WrapWithRoles(customerHandler.GetCustomerByIDHandler, constants.RoleAdmin, constants.RoleManager, constants.RoleEmployee))
	router.GET("/customers", httputil.WrapWithRoles(customerHandler.ListCustomersHandler, constants.RoleAdmin, constants.RoleManager, constants.RoleEmployee))
	router.POST("/customers", httputil.WrapWithRoles(customerHandler.CreateCustomerHandler, constants.RoleAdmin, constants.RoleManager, constants.RoleEmployee))
	router.PUT("/customers/:id", httputil.WrapWithRoles(customerHandler.UpdateCustomerHandler, constants.RoleAdmin, constants.RoleManager))
	router.DELETE("/customers/:id", httputil.WrapWithRole(customerHandler.DeleteCustomerHandler, constants.RoleAdmin))

	// Orders routes - More granular control
	router.GET("/orders", httputil.WrapWithRoles(orderHandler.GetAllOrdersHandler, constants.RoleAdmin, constants.RoleManager, constants.RoleEmployee))
	router.GET("/orders/:id", httputil.Wrap(orderHandler.GetOrderByIDHandler, "read:orders"))
	router.POST("/orders", httputil.Wrap(orderHandler.CreateOrderHandler, "write:orders"))
	router.PUT("/orders/:id", httputil.WrapWithOwnerOrAdmin(orderHandler.UpdateOrderHandler, httputil.ExtractOrderOwnerID))
	router.DELETE("/orders/:id", httputil.WrapWithOwnerOrAdmin(orderHandler.DeleteOrderHandler, httputil.ExtractOrderOwnerID))

	// Sales reports - Admin and Manager only
	router.GET("/sales-reports", httputil.WrapWithRoles(GetSalesReportHandler, constants.RoleAdmin, constants.RoleManager))

	return router
}
