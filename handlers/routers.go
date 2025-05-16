package handlers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/constants"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/inmemorystores"
	"um6p.ma/finalproject/middlewares"
	"um6p.ma/finalproject/models"
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

	// Public routes (no authentication required)
	router.POST("/login", LoginUser)
	router.POST("/register", RegisterUser)

	// Books routes
	router.GET("/books/:id", wrapMiddleware(bookHandler.GetBookByIDHandler, "read:books"))
	router.GET("/books", wrapMiddleware(bookHandler.SearchBooksHandler, "read:books"))
	router.POST("/books", wrapRolesMiddleware(bookHandler.CreateBookHandler, constants.RoleAdmin, constants.RoleManager))
	router.PUT("/books/:id", wrapRolesMiddleware(bookHandler.UpdateBookHandler, constants.RoleAdmin, constants.RoleManager))
	router.DELETE("/books/:id", wrapRoleMiddleware(bookHandler.DeleteBookHandler, constants.RoleAdmin))

	// Authors routes
	router.GET("/authors/:id", wrapMiddleware(authorHandler.GetAuthorByIDHandler, "read:authors"))
	router.GET("/authors", wrapMiddleware(authorHandler.ListAuthorsHandler, "read:authors"))
	router.POST("/authors", wrapRolesMiddleware(authorHandler.CreateAuthorHandler, constants.RoleAdmin, constants.RoleManager))
	router.PUT("/authors/:id", wrapRolesMiddleware(authorHandler.UpdateAuthorHandler, constants.RoleAdmin, constants.RoleManager))
	router.DELETE("/authors/:id", wrapRoleMiddleware(authorHandler.DeleteAuthorHandler, constants.RoleAdmin))

	// Customers routes
	router.GET("/customers/:id", wrapRolesMiddleware(customerHandler.GetCustomerByIDHandler, constants.RoleAdmin, constants.RoleManager, constants.RoleEmployee))
	router.GET("/customers", wrapRolesMiddleware(customerHandler.ListCustomersHandler, constants.RoleAdmin, constants.RoleManager, constants.RoleEmployee))
	router.POST("/customers", wrapRolesMiddleware(customerHandler.CreateCustomerHandler, constants.RoleAdmin, constants.RoleManager, constants.RoleEmployee))
	router.PUT("/customers/:id", wrapRolesMiddleware(customerHandler.UpdateCustomerHandler, constants.RoleAdmin, constants.RoleManager))
	router.DELETE("/customers/:id", wrapRoleMiddleware(customerHandler.DeleteCustomerHandler, constants.RoleAdmin))

	// Orders routes - More granular control
	router.GET("/orders", wrapRolesMiddleware(orderHandler.GetAllOrdersHandler, constants.RoleAdmin, constants.RoleManager, constants.RoleEmployee))
	router.GET("/orders/:id", wrapMiddleware(orderHandler.GetOrderByIDHandler, "read:orders"))
	router.POST("/orders", wrapMiddleware(orderHandler.CreateOrderHandler, "write:orders"))
	router.PUT("/orders/:id", wrapOwnerOrAdminMiddleware(orderHandler.UpdateOrderHandler, extractOrderOwnerID))
	router.DELETE("/orders/:id", wrapOwnerOrAdminMiddleware(orderHandler.DeleteOrderHandler, extractOrderOwnerID))

	// Sales reports - Admin and Manager only
	router.GET("/sales-reports", wrapRolesMiddleware(GetSalesReportHandler, constants.RoleAdmin, constants.RoleManager))

	return router
}

// Middleware wrapper functions
func wrapMiddleware(handler httprouter.Handle, permission string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Chain middleware: Auth -> Permission -> Handler
		finalHandler := middlewares.RequirePermission(permission)(
			middlewares.AuthMiddleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handler(w, r, ps)
				}),
			),
		)
		finalHandler.ServeHTTP(w, r)
	}
}

func wrapRoleMiddleware(handler httprouter.Handle, role string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		finalHandler := middlewares.RequireRole(role)(
			middlewares.AuthMiddleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handler(w, r, ps)
				}),
			),
		)
		finalHandler.ServeHTTP(w, r)
	}
}

func wrapRolesMiddleware(handler httprouter.Handle, roles ...string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		finalHandler := middlewares.RequireRole(roles...)(
			middlewares.AuthMiddleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handler(w, r, ps)
				}),
			),
		)
		finalHandler.ServeHTTP(w, r)
	}
}

func wrapOwnerOrAdminMiddleware(handler httprouter.Handle, extractOwnerID func(*http.Request) int) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		finalHandler := middlewares.OwnerOrAdminMiddleware(extractOwnerID,
			middlewares.AuthMiddleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handler(w, r, ps)
				}),
			),
		)
		finalHandler.ServeHTTP(w, r)
	}
}

func extractOrderOwnerID(r *http.Request) int {
	params := httprouter.ParamsFromContext(r.Context())
	orderID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		return 0
	}

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return 0
	}
	return order.CustomerID
}
