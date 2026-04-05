package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	var create model.OrderCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	order, err := h.services.Order.CreateOrder(r.Context(), userID, create)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	Success(w, order)
}

func (h *Handler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	orders, err := h.services.Order.GetUserOrders(r.Context(), userID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, orders)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	orderIDStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		BadRequest(w, "Invalid order ID")
		return
	}

	order, err := h.services.Order.GetOrder(r.Context(), orderID, requester, role)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			NotFound(w, "Order not found")
			return
		}
		NotFound(w, err.Error())
		return
	}
	Success(w, order)
}

func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	orderIDStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		BadRequest(w, "Invalid order ID")
		return
	}

	var update struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if err := h.services.Order.UpdateOrderStatus(r.Context(), orderID, update.Status, requester, role); err != nil {
		switch {
		case errors.Is(err, apperr.ErrForbidden):
			Forbidden(w, "Not allowed to update this order")
		default:
			BadRequest(w, err.Error())
		}
		return
	}

	order, err := h.services.Order.GetOrder(r.Context(), orderID, requester, role)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			NotFound(w, "Order not found")
			return
		}
		InternalServerError(w, err.Error())
		return
	}
	Success(w, order)
}
