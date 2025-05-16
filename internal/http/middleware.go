package http

import (
	"github.com/julienschmidt/httprouter"
)

// Wrap wraps a handler with authentication and permission middleware
func Wrap(handler httprouter.Handle, permission string) httprouter.Handle {
	return WrapWithMiddleware(handler,
		RequireAuth,
		RequirePermission(permission),
	)
}

// WrapWithRole wraps a handler with authentication and single role middleware
func WrapWithRole(handler httprouter.Handle, role string) httprouter.Handle {
	return WrapWithMiddleware(handler,
		RequireAuth,
		RequireRoles(role),
	)
}

// WrapWithRoles wraps a handler with authentication and multiple roles middleware
func WrapWithRoles(handler httprouter.Handle, roles ...string) httprouter.Handle {
	return WrapWithMiddleware(handler,
		RequireAuth,
		RequireRoles(roles...),
	)
}

// WrapWithOwnerOrAdmin wraps a handler with authentication and owner/admin check middleware
func WrapWithOwnerOrAdmin(handler httprouter.Handle, extractOwnerID ExtractOwnerIDFunc) httprouter.Handle {
	return WrapWithMiddleware(handler,
		RequireAuth,
		RequireOwnerOrAdmin(extractOwnerID),
	)
}
