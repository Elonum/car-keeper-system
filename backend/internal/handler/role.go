package handler

import (
	"net/http"

	"github.com/carkeeper/backend/internal/authz"
)

// AdminListRoleDefinitions returns all roles (for admin UI and integration); admin permission only.
func (h *Handler) AdminListRoleDefinitions(w http.ResponseWriter, r *http.Request) {
	if _, ok := RequirePermission(w, r, authz.PermAdminRolesView); !ok {
		return
	}
	list, err := h.services.Role.ListDefinitions(r.Context())
	if err != nil {
		HandleError(w, r, err)
		return
	}
	Success(w, list)
}
