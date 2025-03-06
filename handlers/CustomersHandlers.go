package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/interfaces"
	"um6p.ma/finalproject/models"
)

type CustomerHandler struct {
	Store interfaces.CustomerStore
}

// GetCustomerByIDHandler retrieves a customer by ID from the database
func (h *CustomerHandler) GetCustomerByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	if err := database.DB.WithContext(ctx).First(&customer, id).Error; err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert into the database
	if err := database.DB.WithContext(ctx).Create(&newCustomer).Error; err != nil {
		http.Error(w, "Failed to add customer", http.StatusInternalServerError)
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
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedCustomer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&updatedCustomer); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update in the database
	if err := database.DB.WithContext(ctx).Model(&models.Customer{}).Where("id = ?", id).Updates(updatedCustomer).Error; err != nil {
		http.Error(w, "Customer not found or update failed", http.StatusNotFound)
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
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete from the database
	if err := database.DB.WithContext(ctx).Delete(&models.Customer{}, id).Error; err != nil {
		http.Error(w, "Customer not found or delete failed", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListCustomersHandler retrieves all customers from the database
func (h *CustomerHandler) ListCustomersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var customers []models.Customer

	if err := database.DB.WithContext(ctx).Find(&customers).Error; err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}
