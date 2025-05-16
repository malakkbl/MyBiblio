package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"um6p.ma/finalproject/constants"
)

// Secret key for JWT
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Context Key for storing user data
type UserContextKey string

const (
	ContextUserKey UserContextKey = "user"
)

// Claims represents the structure of our custom JWT claims
type Claims struct {
	UserID      int      `json:"user_id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

// AuthMiddleware verifies the JWT token and adds user info to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Parse JWT Token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), ContextUserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware checks if user has the required role
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ContextUserKey).(*Claims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Admin has access to everything
			if claims.Role == constants.RoleAdmin {
				next.ServeHTTP(w, r)
				return
			}

			// Check if user's role matches any of the required roles
			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden: Insufficient role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission middleware checks if user has the required permission
func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ContextUserKey).(*Claims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has the required permission
			hasPermission := false
			for _, p := range claims.Permissions {
				if p == permission || p == "write:all" || p == "read:all" {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext extracts user claims from the request context
func GetUserFromContext(r *http.Request) (*Claims, bool) {
	claims, ok := r.Context().Value(ContextUserKey).(*Claims)
	return claims, ok
}

// AdminOnlyMiddleware is a shortcut for requiring admin role
func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return RequireRole(constants.RoleAdmin)(next)
}

// OwnerOrAdminMiddleware allows access only to the resource owner or an admin
func OwnerOrAdminMiddleware(extractOwnerID func(*http.Request) int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetUserFromContext(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Admins can proceed without checks
		if claims.Role == constants.RoleAdmin {
			next.ServeHTTP(w, r)
			return
		}

		// Check if the user owns the resource
		resourceOwnerID := extractOwnerID(r)
		if resourceOwnerID != claims.UserID {
			http.Error(w, "Forbidden: You are not the owner", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
