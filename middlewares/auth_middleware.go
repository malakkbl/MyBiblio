package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key for JWT
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Context Key for storing user data
type UserContextKey string

const ContextUserKey UserContextKey = "user"

// AuthMiddleware verifies the JWT token
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

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ExtractUserInfo extracts the user ID and role from the request context
func ExtractUserInfo(r *http.Request) (int, string, bool) {
	claims, ok := r.Context().Value(ContextUserKey).(jwt.MapClaims)
	if !ok {
		return 0, "", false
	}

	userID, ok1 := claims["user_id"].(float64) // JWT stores numbers as float64
	role, ok2 := claims["role"].(string)
	if !ok1 || !ok2 {
		return 0, "", false
	}

	return int(userID), role, true
}

// AdminOnlyMiddleware restricts access to admins
func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, role, ok := ExtractUserInfo(r)
		if !ok || role != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// OwnerOrAdminMiddleware allows access only to the resource owner or an admin
func OwnerOrAdminMiddleware(extractOwnerID func(*http.Request) int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, role, ok := ExtractUserInfo(r)
		if !ok {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Admins can proceed without checks
		if role == "admin" {
			next.ServeHTTP(w, r)
			return
		}

		// Check if the user owns the resource
		resourceOwnerID := extractOwnerID(r)
		if resourceOwnerID != userID {
			http.Error(w, "Forbidden: You are not the owner", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
