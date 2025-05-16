package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/middlewares"
	"um6p.ma/finalproject/models"
)

// WrapMiddleware wraps a handler with authentication and permission middleware
func WrapMiddleware(handler httprouter.Handle, permission string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Chain middleware: Auth -> Permission -> Handler
		httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, ps)
		})

		// Chain middlewares
		finalHandler := middlewares.RequirePermission(permission)(
			middlewares.AuthMiddleware(httpHandler),
		)

		// Handle any panics that might occur
		defer func() {
			if err := recover(); err != nil {
				errorhandling.HandleError(w, errorhandling.NewError(
					http.StatusInternalServerError,
					errorhandling.ErrCodeInternalServer,
					"Internal server error",
				).WithDebug(fmt.Sprint(err)))
			}
		}()

		finalHandler.ServeHTTP(w, r)
	}
}

// WrapRoleMiddleware wraps a handler with authentication and single role middleware
func WrapRoleMiddleware(handler httprouter.Handle, role string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, ps)
		})

		finalHandler := middlewares.RequireRole(role)(
			middlewares.AuthMiddleware(httpHandler),
		)

		defer func() {
			if err := recover(); err != nil {
				errorhandling.HandleError(w, errorhandling.NewError(
					http.StatusInternalServerError,
					errorhandling.ErrCodeInternalServer,
					"Internal server error",
				).WithDebug(fmt.Sprint(err)))
			}
		}()

		finalHandler.ServeHTTP(w, r)
	}
}

// WrapRolesMiddleware wraps a handler with authentication and multiple roles middleware
func WrapRolesMiddleware(handler httprouter.Handle, roles ...string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, ps)
		})

		finalHandler := middlewares.RequireRole(roles...)(
			middlewares.AuthMiddleware(httpHandler),
		)

		defer func() {
			if err := recover(); err != nil {
				errorhandling.HandleError(w, errorhandling.NewError(
					http.StatusInternalServerError,
					errorhandling.ErrCodeInternalServer,
					"Internal server error",
				).WithDebug(fmt.Sprint(err)))
			}
		}()

		finalHandler.ServeHTTP(w, r)
	}
}

// WrapOwnerOrAdminMiddleware wraps a handler with authentication and owner/admin check middleware
func WrapOwnerOrAdminMiddleware(handler httprouter.Handle, extractOwnerID func(*http.Request) int) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, ps)
		})

		finalHandler := middlewares.OwnerOrAdminMiddleware(extractOwnerID,
			middlewares.AuthMiddleware(httpHandler),
		)

		defer func() {
			if err := recover(); err != nil {
				errorhandling.HandleError(w, errorhandling.NewError(
					http.StatusInternalServerError,
					errorhandling.ErrCodeInternalServer,
					"Internal server error",
				).WithDebug(fmt.Sprint(err)))
			}
		}()

		finalHandler.ServeHTTP(w, r)
	}
}

// ExtractOrderOwnerID extracts the owner ID from an order
func ExtractOrderOwnerID(r *http.Request) int {
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
