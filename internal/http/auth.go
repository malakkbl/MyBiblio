package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"um6p.ma/finalproject/errorhandling"
)

// ContextKey is a type for context keys
type ContextKey string

// ExtractOwnerIDFunc is a function type for extracting resource owner IDs
type ExtractOwnerIDFunc func(*http.Request) int

// Context keys
const (
	ContextUserKey ContextKey = "user"
)

// Claims represents JWT claims
type Claims struct {
	UserID      int      `json:"user_id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

// GetJWTSecret returns the JWT secret key used for signing and validating tokens
func GetJWTSecret() []byte {
	return jwtSecret
}

// RequireAuth is middleware that validates JWT tokens
var RequireAuth MiddlewareFunc = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorhandling.HandleError(w, errorhandling.ErrMissingToken)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			errorhandling.HandleError(w, errorhandling.ErrInvalidToken.
				WithDebug("Invalid token format"))
			return
		}

		// Parse JWT Token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				errorhandling.HandleError(w, errorhandling.ErrInvalidToken)
				return
			}
			errorhandling.HandleError(w, errorhandling.ErrExpiredToken)
			return
		}

		if !token.Valid {
			errorhandling.HandleError(w, errorhandling.ErrInvalidToken)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), ContextUserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireOwnerOrAdmin creates middleware that ensures the user is either the resource owner or an admin
func RequireOwnerOrAdmin(extractOwnerID ExtractOwnerIDFunc) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaimsFromContext(r.Context())
			if !ok {
				errorhandling.HandleError(w, errorhandling.ErrMissingToken)
				return
			}

			// Admins can access everything
			if claims.Role == "admin" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if user is the resource owner
			ownerID := extractOwnerID(r)
			if ownerID == 0 || ownerID != claims.UserID {
				errorhandling.HandleError(w, errorhandling.NewError(
					http.StatusForbidden,
					errorhandling.ErrCodeForbidden,
					"Access denied: you are not the owner of this resource",
				))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
