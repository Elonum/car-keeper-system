package service

import (
	"github.com/carkeeper/backend/internal/authz"
	"github.com/carkeeper/backend/internal/model"
)

// UserResponseFrom builds a client-safe profile including RBAC from the live authz matrix.
func UserResponseFrom(u *model.User) model.UserResponse {
	resp := u.ToResponse()
	resp.Permissions = authz.PermissionsForRole(u.Role)
	resp.IsStaff = authz.IsStaff(u.Role)
	return resp
}
