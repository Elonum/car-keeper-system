package handler

import (
	"net/http"

	"github.com/carkeeper/backend/internal/authz"
	"github.com/carkeeper/backend/internal/middleware"
	"github.com/google/uuid"
)

// RequesterAndRole returns the authenticated user and role, or writes 401 and ok=false.
func RequesterAndRole(w http.ResponseWriter, r *http.Request) (requester uuid.UUID, role string, ok bool) {
	requester, ok = middleware.UserIDFromContext(r.Context())
	if !ok {
		Unauthorized(w, "User not authenticated")
		return uuid.Nil, "", false
	}
	role, _ = middleware.GetUserRole(r.Context())
	return requester, role, true
}

// RequireStaff requires an authenticated admin or manager.
func RequireStaff(w http.ResponseWriter, r *http.Request) (requester uuid.UUID, ok bool) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return uuid.Nil, false
	}
	if !authz.IsStaff(role) {
		Forbidden(w, "Only admins and managers can perform this action")
		return uuid.Nil, false
	}
	return requester, true
}
