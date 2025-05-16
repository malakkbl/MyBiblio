package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"um6p.ma/finalproject/constants"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/models"
)

// Secret key for JWT (from environment variable)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// RegisterUser handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate role
	if !constants.ValidRoles[input.Role] {
		http.Error(w, "Invalid role. Must be one of: admin, user, manager, employee", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var existingUser models.User
	result := database.DB.Where("email = ?", input.Email).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create new user object
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	// Save user in the database
	if err := database.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// LoginUser handles user authentication and token generation
func LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var user models.User

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user in database
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT Token with enhanced claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     user.ID,
		"email":       user.Email,
		"name":        user.Name,
		"role":        user.Role,
		"permissions": getPermissionsForRole(user.Role),
		"iat":         time.Now().Unix(),
		"exp":         time.Now().Add(time.Hour * 24).Unix(), // Expires in 24h
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Return JWT Token with user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// getPermissionsForRole returns a list of permissions based on the user's role
func getPermissionsForRole(role string) []string {
	switch role {
	case constants.RoleAdmin:
		return []string{
			"read:all",
			"write:all",
			"delete:all",
			"manage:users",
			"manage:roles",
			"generate:reports",
		}
	case constants.RoleManager:
		return []string{
			"read:all",
			"write:books",
			"write:authors",
			"write:orders",
			"generate:reports",
		}
	case constants.RoleEmployee:
		return []string{
			"read:all",
			"write:orders",
			"write:customers",
		}
	case constants.RoleUser:
		return []string{
			"read:books",
			"read:authors",
			"write:orders",
		}
	default:
		return []string{"read:books"}
	}
}
