package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/interfaces"
	"um6p.ma/finalproject/models"
	"um6p.ma/finalproject/validation"
)

type CustomerHandler struct {
	Store interfaces.CustomerStore
}

// GetCustomerByIDHandler retrieves a customer by ID from the database
func (h *CustomerHandler) GetCustomerByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid ID format",
		))
		return
	}

	var customer models.Customer
	if err := database.DB.WithContext(ctx).First(&customer, id).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewNotFoundError("Customer", id))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// CreateCustomerHandler adds a new customer to the database
func (h *CustomerHandler) CreateCustomerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var newCustomer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid request body",
		).WithDebug(err.Error()))
		return
	}

	// Validate customer data
	if errors := validation.Validate(newCustomer); len(errors) > 0 {
		errorhandling.HandleError(w, errorhandling.NewValidationError(errors))
		return
	}

	// Check if email is unique
	var existingCustomer models.Customer
	if err := database.DB.WithContext(ctx).Where("email = ?", newCustomer.Email).First(&existingCustomer).Error; err == nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusConflict,
			errorhandling.ErrCodeDuplicateEntry,
			"Email already registered",
		))
		return
	}

	// Insert into the database
	if err := database.DB.WithContext(ctx).Create(&newCustomer).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCustomer)
}

// UpdateCustomerHandler modifies an existing customer in the database
func (h *CustomerHandler) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid ID format",
		))
		return
	}

	var updatedCustomer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&updatedCustomer); err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid request body",
		).WithDebug(err.Error()))
		return
	}

	// Validate customer data
	if errors := validation.Validate(updatedCustomer); len(errors) > 0 {
		errorhandling.HandleError(w, errorhandling.NewValidationError(errors))
		return
	}

	// Check if customer exists
	var existingCustomer models.Customer
	if err := database.DB.WithContext(ctx).First(&existingCustomer, id).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewNotFoundError("Customer", id))
		return
	}

	// Check if new email is unique (if changed)
	if updatedCustomer.Email != existingCustomer.Email {
		var emailExists models.Customer
		if err := database.DB.WithContext(ctx).Where("email = ? AND id != ?", updatedCustomer.Email, id).First(&emailExists).Error; err == nil {
			errorhandling.HandleError(w, errorhandling.NewError(
				http.StatusConflict,
				errorhandling.ErrCodeDuplicateEntry,
				"Email already registered",
			))
			return
		}
	}

	// Update in the database
	if err := database.DB.WithContext(ctx).Model(&existingCustomer).Updates(updatedCustomer).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCustomer)
}

// DeleteCustomerHandler removes a customer from the database
func (h *CustomerHandler) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid ID format",
		))
		return
	}

	// Check if customer has any orders
	var orderCount int64
	if err := database.DB.WithContext(ctx).Model(&models.Order{}).Where("customer_id = ?", id).Count(&orderCount).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	if orderCount > 0 {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusConflict,
			errorhandling.ErrCodeBadRequest,
			"Cannot delete customer with existing orders",
		).WithDetails(map[string]interface{}{
			"orderCount": orderCount,
		}))
		return
	}

	// Delete from the database
	if err := database.DB.WithContext(ctx).Delete(&models.Customer{}, id).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListCustomersHandler retrieves all customers from the database
func (h *CustomerHandler) ListCustomersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var customers []models.Customer

	query := database.DB.WithContext(ctx)

	// Handle query parameters for filtering
	if email := r.URL.Query().Get("email"); email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if name := r.URL.Query().Get("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if city := r.URL.Query().Get("city"); city != "" {
		query = query.Where("address_city LIKE ?", "%"+city+"%")
	}

	if err := query.Find(&customers).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}
