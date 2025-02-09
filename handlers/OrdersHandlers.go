package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/inmemorystores"
	"um6p.ma/finalproject/interfaces"
	"um6p.ma/finalproject/models"
)

type OrderHandler struct {
	Store       interfaces.OrderStore
	ReportStore *inmemorystores.ReportStore
}

func (h *OrderHandler) GetOrderByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	order, err := h.Store.GetOrder(ctx, id)
	if err != nil {
		if err == errorhandling.ErrOrderNotFound {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdOrder, err := h.Store.CreateOrder(ctx, newOrder)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdOrder)
}

func (h *OrderHandler) UpdateOrderHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedOrder models.Order
	err = json.NewDecoder(r.Body).Decode(&updatedOrder)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.Store.UpdateOrder(ctx, id, updatedOrder)
	if err != nil {
		if err == errorhandling.ErrOrderNotFound {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.Store.DeleteOrder(ctx, id)
	if err != nil {
		if err == errorhandling.ErrOrderNotFound {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) GetAllOrdersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	orders, err := h.Store.GetAllOrders(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetSalesReportHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	report := h.ReportStore.GetLatestReport()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
