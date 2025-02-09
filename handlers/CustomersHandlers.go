package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/interfaces"
	"um6p.ma/finalproject/models"
)

type CustomerHandler struct {
	Store interfaces.CustomerStore
}

func (h *CustomerHandler) GetCustomerByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	customer, err := h.Store.GetCustomer(ctx, id)
	if err != nil {
		if err == errorhandling.ErrCustomerNotFound {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) CreateCustomerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var newCustomer models.Customer
	err := json.NewDecoder(r.Body).Decode(&newCustomer)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdCustomer, err := h.Store.CreateCustomer(ctx, newCustomer)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCustomer)
}

func (h *CustomerHandler) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedCustomer models.Customer
	err = json.NewDecoder(r.Body).Decode(&updatedCustomer)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	customer, err := h.Store.UpdateCustomer(ctx, id, updatedCustomer)
	if err != nil {
		if err == errorhandling.ErrCustomerNotFound {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.Store.DeleteCustomer(ctx, id)
	if err != nil {
		if err == errorhandling.ErrCustomerNotFound {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CustomerHandler) ListCustomersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	customers, err := h.Store.GetAllCustomers(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}
