package handler

import (
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
	if !DecodeJSON(w, r, &create) {
		return
	}

	order, err := h.services.Order.CreateOrder(r.Context(), userID, create)
	if err != nil {
		HandleError(w, r, err)
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
		HandleError(w, r, err)
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
		HandleError(w, r, err)
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
	if !DecodeJSON(w, r, &update) {
		return
	}

	if err := h.services.Order.UpdateOrderStatus(r.Context(), orderID, update.Status, requester, role); err != nil {
		if errors.Is(err, apperr.ErrForbidden) {
			Forbidden(w, "Not allowed to update this order")
			return
		}
		HandleError(w, r, err)
		return
	}

	order, err := h.services.Order.GetOrder(r.Context(), orderID, requester, role)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, order)
}
