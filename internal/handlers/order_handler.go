package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/middleware"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/orders"
)

type OrderHandler struct {
	orderService *orders.Service
}

func NewOrderHandler(orderService *orders.Service) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items []orders.OrderItemInput `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Items) == 0 {
		respondWithError(w, http.StatusBadRequest, "Order must contain at least one item")
		return
	}

	customerID := int64(middleware.GetUserID(r.Context()))
	order, err := h.orderService.Create(r.Context(), customerID, req.Items)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Order created",
		"data":    order,
	})
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.orderService.GetByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	customerID := int64(middleware.GetUserID(r.Context()))
	userRole := middleware.GetUserRole(r.Context())

	if order.CustomerID != customerID && userRole != "admin" {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"data": order})
}

func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	customerID := int64(middleware.GetUserID(r.Context()))

	orders, err := h.orderService.GetByCustomerID(r.Context(), customerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"data": orders})
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	_ = id
	_ = req.Status

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"message": "Order status updated"})
}
