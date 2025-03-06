package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/inmemorystores"
	"um6p.ma/finalproject/interfaces"
	"um6p.ma/finalproject/models"
)

type OrderHandler struct {
	Store       interfaces.OrderStore
	ReportStore *inmemorystores.ReportStore
}

// GetOrderByIDHandler retrieves an order by ID from the database
func (h *OrderHandler) GetOrderByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var order models.Order
	if err := database.DB.WithContext(ctx).Preload("Items").First(&order, id).Error; err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// CreateOrderHandler adds a new order to the database
func (h *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var newOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert into the database
	if err := database.DB.WithContext(ctx).Create(&newOrder).Error; err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
}

// UpdateOrderHandler modifies an existing order in the database
func (h *OrderHandler) UpdateOrderHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update in the database
	if err := database.DB.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Updates(updatedOrder).Error; err != nil {
		http.Error(w, "Order not found or update failed", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

// DeleteOrderHandler removes an order from the database
func (h *OrderHandler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete from the database
	if err := database.DB.WithContext(ctx).Delete(&models.Order{}, id).Error; err != nil {
		http.Error(w, "Order not found or delete failed", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllOrdersHandler retrieves all orders from the database
func (h *OrderHandler) GetAllOrdersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var orders []models.Order

	if err := database.DB.WithContext(ctx).Preload("Items").Find(&orders).Error; err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetSalesReportHandler returns the latest sales reports
func (h *OrderHandler) GetSalesReportHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var reports []models.SalesReport

	if err := database.DB.WithContext(ctx).Order("timestamp DESC").Limit(10).Find(&reports).Error; err != nil {
		http.Error(w, "Failed to fetch sales reports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
