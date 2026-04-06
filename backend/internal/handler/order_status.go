package handler

import (
	"net/http"

	"github.com/carkeeper/backend/internal/authz"
	"github.com/carkeeper/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetOrderStatuses(w http.ResponseWriter, r *http.Request) {
	list, err := h.services.OrderStatus.ListPublic(r.Context())
	if err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, list)
}

func (h *Handler) AdminListOrderStatuses(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermOrdersManageStatus); !ok {
		return
	}
	list, err := h.services.OrderStatus.ListAll(r.Context())
	if err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, list)
}

func (h *Handler) AdminCreateOrderStatus(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermAdminOrderStatuses); !ok {
		return
	}
	var in model.OrderStatusCreate
	if !DecodeJSON(w, r, &in) {
		return
	}
	d, err := h.services.OrderStatus.Create(r.Context(), in)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	JSON(w, http.StatusCreated, Response{Success: true, Data: d})
}

func (h *Handler) AdminUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermAdminOrderStatuses); !ok {
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		BadRequest(w, "Invalid status ID")
		return
	}
	var patch model.OrderStatusUpdate
	if !DecodeJSON(w, r, &patch) {
		return
	}
	d, err := h.services.OrderStatus.Update(r.Context(), id, patch)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, d)
}

func (h *Handler) AdminDeleteOrderStatus(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermAdminOrderStatuses); !ok {
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		BadRequest(w, "Invalid status ID")
		return
	}
	if err := h.services.OrderStatus.Delete(r.Context(), id); err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, map[string]string{"message": "Order status deleted"})
}
