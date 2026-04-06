package handler

import (
	"net/http"

	"github.com/carkeeper/backend/internal/authz"
)

// AdminListAllOrders returns all orders (staff with orders.view_any).
func (h *Handler) AdminListAllOrders(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermOrdersViewAny); !ok {
		return
	}
	list, err := h.services.Order.ListAllOrdersForStaff(r.Context())
	if err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, list)
}

// AdminListAllAppointments returns all service appointments (staff with appointments.view_any).
func (h *Handler) AdminListAllAppointments(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermAppointmentsViewAny); !ok {
		return
	}
	list, err := h.services.Service.ListAllAppointmentsForStaff(r.Context())
	if err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, list)
}
