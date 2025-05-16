package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/models"
)

// MiddlewareFunc is a type alias for HTTP middleware
type MiddlewareFunc func(http.Handler) http.Handler

// WrapWithMiddleware wraps a handler with multiple middleware functions
func WrapWithMiddleware(handler httprouter.Handle, mw ...MiddlewareFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, ps)
		})
		// Chain middleware in reverse order
		for i := len(mw) - 1; i >= 0; i-- {
			h = http.HandlerFunc(mw[i](h).ServeHTTP)
		}

		defer func() {
			if err := recover(); err != nil {
				errorhandling.HandleError(w, errorhandling.NewError(
					http.StatusInternalServerError,
					errorhandling.ErrCodeInternalServer,
					"Internal server error",
				).WithDebug(fmt.Sprint(err)))
			}
		}()

		h.ServeHTTP(w, r)
	}
}

// RequirePermission creates a middleware that checks for a specific permission
func RequirePermission(permission string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaimsFromContext(r.Context())
			if !ok {
				errorhandling.HandleError(w, errorhandling.ErrMissingToken)
				return
			}

			if !HasPermission(claims.Permissions, permission) {
				errorhandling.HandleError(w, errorhandling.NewError(
					http.StatusForbidden,
					errorhandling.ErrCodeForbidden,
					"Insufficient permissions",
				))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRoles creates a middleware that checks for specific roles
func RequireRoles(roles ...string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaimsFromContext(r.Context())
			if !ok {
				errorhandling.HandleError(w, errorhandling.ErrMissingToken)
				return
			}

			if !HasRole(claims.Role, roles...) {
				errorhandling.HandleError(w, errorhandling.NewError(
					http.StatusForbidden,
					errorhandling.ErrCodeForbidden,
					fmt.Sprintf("Access denied. Required roles: %v", roles),
				))
				return
			}

			next.ServeHTTP(w, r)
		})
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

// GetClaimsFromContext extracts JWT claims from the request context
func GetClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(ContextUserKey).(*Claims)
	return claims, ok
}

// HasPermission checks if the given permissions include the required one
func HasPermission(permissions []string, required string) bool {
	for _, p := range permissions {
		if p == required || p == "*" {
			return true
		}
	}
	return false
}

// HasRole checks if the given role matches any of the required roles
func HasRole(userRole string, requiredRoles ...string) bool {
	for _, role := range requiredRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
