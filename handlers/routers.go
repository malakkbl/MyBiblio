package handlers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/inmemorystores"
	"um6p.ma/finalproject/middlewares"
	"um6p.ma/finalproject/models"
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

	// Authentication (No protection)
	router.POST("/login", LoginUser)
	router.POST("/register", RegisterUser)

	// Protected Routes - Authenticated Users (Any user can access)
	router.GET("/books/:id", wrapMiddleware(bookHandler.GetBookByIDHandler))
	router.GET("/books", wrapMiddleware(bookHandler.SearchBooksHandler))
	router.GET("/authors/:id", wrapMiddleware(authorHandler.GetAuthorByIDHandler))
	router.GET("/authors", wrapMiddleware(authorHandler.ListAuthorsHandler))
	router.GET("/customers/:id", wrapMiddleware(customerHandler.GetCustomerByIDHandler))
	router.GET("/customers", wrapMiddleware(customerHandler.ListCustomersHandler))
	router.GET("/orders", wrapMiddleware(orderHandler.GetAllOrdersHandler))
	router.GET("/orders/:id", wrapMiddleware(orderHandler.GetOrderByIDHandler))
	router.GET("/sales-reports", wrapMiddleware(GetSalesReportHandler))

	// Admin Only Routes
	router.POST("/books", wrapAdminMiddleware(bookHandler.CreateBookHandler))
	router.PUT("/books/:id", wrapAdminMiddleware(bookHandler.UpdateBookHandler))
	router.DELETE("/books/:id", wrapAdminMiddleware(bookHandler.DeleteBookHandler))

	router.POST("/authors", wrapAdminMiddleware(authorHandler.CreateAuthorHandler))
	router.PUT("/authors/:id", wrapAdminMiddleware(authorHandler.UpdateAuthorHandler))
	router.DELETE("/authors/:id", wrapAdminMiddleware(authorHandler.DeleteAuthorHandler))

	router.POST("/customers", wrapAdminMiddleware(customerHandler.CreateCustomerHandler))
	router.PUT("/customers/:id", wrapAdminMiddleware(customerHandler.UpdateCustomerHandler))
	router.DELETE("/customers/:id", wrapAdminMiddleware(customerHandler.DeleteCustomerHandler))

	// Orders - Any user can create an order, but only the owner or admin can modify/delete
	router.POST("/orders", wrapMiddleware(orderHandler.CreateOrderHandler))
	router.PUT("/orders/:id", wrapOwnerOrAdminMiddleware(orderHandler.UpdateOrderHandler, extractOrderOwnerID))
	router.DELETE("/orders/:id", wrapOwnerOrAdminMiddleware(orderHandler.DeleteOrderHandler, extractOrderOwnerID))

	return router
}

// Middleware Wrappers
func wrapMiddleware(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Apply authentication middleware
		authMiddleware := middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, ps)
		}))

		// Serve request
		authMiddleware.ServeHTTP(w, r)
	}
}

func wrapAdminMiddleware(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Apply authentication and admin check middleware
		adminMiddleware := middlewares.AdminOnlyMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, ps)
		})))

		// Serve request
		adminMiddleware.ServeHTTP(w, r)
	}
}

func wrapOwnerOrAdminMiddleware(handler httprouter.Handle, extractOwnerID func(*http.Request) int) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		middleware := middlewares.OwnerOrAdminMiddleware(extractOwnerID, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, ps)
		}))
		middleware.ServeHTTP(w, r)
	}
}

func extractOrderOwnerID(r *http.Request) int {
	params := httprouter.ParamsFromContext(r.Context())
	orderID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		return 0
	}

	// Fetch the order from the database
	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return 0
	}
	return order.CustomerID
}
