package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"um6p.ma/finalproject/constants"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/errorhandling"
	internalhttp "um6p.ma/finalproject/internal/http"
	"um6p.ma/finalproject/models"
	"um6p.ma/finalproject/validation"
)

type RegisterInput struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100,passwd"`
	Role     string `json:"role" validate:"required,oneof=admin manager employee user"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterUser handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input RegisterInput

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid request format",
		).WithDebug(err.Error()))
		return
	}

	// Validate input
	if errs := validation.Validate(input); len(errs) > 0 {
		errorhandling.HandleError(w, errorhandling.NewValidationError(errs))
		return
	}

	// Check if user already exists
	var existingUser models.User
	result := database.DB.Where("email = ?", input.Email).First(&existingUser)
	if result.Error == nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusConflict,
			errorhandling.ErrCodeDuplicateEntry,
			"Email already registered",
		))
		return
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(result.Error))
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusInternalServerError,
			errorhandling.ErrCodeInternalServer,
			"Failed to process password",
		).WithDebug(err.Error()))
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
			errorhandling.HandleError(w, errorhandling.NewError(
				http.StatusConflict,
				errorhandling.ErrCodeDuplicateEntry,
				"Email already registered",
			))
			return
		}
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
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
	var input LoginInput
	var user models.User

	// Decode and validate request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid request format",
		).WithDebug(err.Error()))
		return
	}

	// Validate input
	if errs := validation.Validate(input); len(errs) > 0 {
		errorhandling.HandleError(w, errorhandling.NewValidationError(errs))
		return
	}

	// Find user in database
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorhandling.HandleError(w, errorhandling.ErrInvalidCredentials)
			return
		}
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		errorhandling.HandleError(w, errorhandling.ErrInvalidCredentials)
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

	tokenString, err := token.SignedString(internalhttp.GetJWTSecret())
	if err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusInternalServerError,
			errorhandling.ErrCodeInternalServer,
			"Failed to generate authentication token",
		).WithDebug(err.Error()))
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
		"permissions": getPermissionsForRole(user.Role),
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
