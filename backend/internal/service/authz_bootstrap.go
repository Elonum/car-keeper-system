package service

import (
	"context"
	"log/slog"

	"github.com/carkeeper/backend/internal/authz"
	"github.com/carkeeper/backend/internal/repository"
)

// BootstrapAuthz loads staff roles and role_permissions from the database into authz.
// On failure, in-memory defaults from authz.DefaultRolePermissions() remain in effect.
func BootstrapAuthz(ctx context.Context, repos *repository.Repository) {
	staff, err := repos.Role.ListStaffCodes(ctx)
	if err != nil {
		slog.Warn("authz bootstrap: staff roles not loaded, using defaults", "err", err)
	} else {
		authz.SetStaffRoles(staff)
		slog.Info("authz bootstrap: staff roles loaded", "count", len(staff))
	}

	perms, err := repos.Role.LoadRolePermissions(ctx)
	if err != nil {
		slog.Warn("authz bootstrap: role permissions not loaded, using defaults", "err", err)
		return
	}
	if len(perms) == 0 {
		slog.Warn("authz bootstrap: role_permissions empty, using default permission matrix")
		return
	}
	authz.SetRolePermissions(perms)
	slog.Info("authz bootstrap: role permissions loaded", "roles", len(perms))
}
