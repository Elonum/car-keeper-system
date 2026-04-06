package authz

import (
	"sync"

	"github.com/google/uuid"
)

// defaultStaffRoles used before SetStaffRoles (tests) or if DB bootstrap fails.
var defaultStaffRoles = []string{"admin", "manager", "service_advisor"}

var (
	staffMu    sync.RWMutex
	staffCodes map[string]struct{} // nil = use defaultStaffRoles in IsStaff

	permMu    sync.RWMutex
	rolePerms map[string]map[string]struct{} // role -> permission set; nil = use DefaultRolePermissions in HasPermission
)

// SetStaffRoles replaces the in-memory staff set (typically from role_definitions.is_staff).
func SetStaffRoles(codes []string) {
	staffMu.Lock()
	defer staffMu.Unlock()
	if len(codes) == 0 {
		staffCodes = nil
		return
	}
	staffCodes = make(map[string]struct{}, len(codes))
	for _, c := range codes {
		staffCodes[c] = struct{}{}
	}
}

// SetRolePermissions replaces the in-memory RBAC matrix from role_permissions.
// Pass nil or empty map to fall back to DefaultRolePermissions() inside HasPermission.
func SetRolePermissions(fromDB map[string][]string) {
	permMu.Lock()
	defer permMu.Unlock()
	if len(fromDB) == 0 {
		rolePerms = nil
		return
	}
	rolePerms = make(map[string]map[string]struct{}, len(fromDB))
	for role, list := range fromDB {
		if len(list) == 0 {
			rolePerms[role] = nil
			continue
		}
		m := make(map[string]struct{}, len(list))
		for _, p := range list {
			m[p] = struct{}{}
		}
		rolePerms[role] = m
	}
}

// HasPermission reports whether role has the given permission_code.
func HasPermission(role, permission string) bool {
	if role == "" || permission == "" {
		return false
	}
	permMu.RLock()
	custom := rolePerms
	permMu.RUnlock()

	var set map[string]struct{}
	if len(custom) > 0 {
		set = custom[role]
	} else {
		def := DefaultRolePermissions()
		list, ok := def[role]
		if !ok {
			return false
		}
		set = make(map[string]struct{}, len(list))
		for _, p := range list {
			set[p] = struct{}{}
		}
	}
	if len(set) == 0 {
		return false
	}
	_, ok := set[permission]
	return ok
}

// IsStaff returns true if the role is marked staff in role_definitions (elevated internal account).
// It does not imply every API permission; use HasPermission for fine-grained checks.
func IsStaff(role string) bool {
	staffMu.RLock()
	set := staffCodes
	staffMu.RUnlock()

	if len(set) > 0 {
		_, ok := set[role]
		return ok
	}
	for _, r := range defaultStaffRoles {
		if r == role {
			return true
		}
	}
	return false
}

// IsAdmin is true only for the administrator role (full platform control).
func IsAdmin(role string) bool {
	return role == "admin"
}

// IsOwnerOrHasPermission is true when requester owns the resource or holds permission.
func IsOwnerOrHasPermission(resourceOwnerID, requester uuid.UUID, role, permission string) bool {
	if resourceOwnerID == requester {
		return true
	}
	return HasPermission(role, permission)
}

// CanViewOrder reports whether the requester may read the order (owner or staff with orders.view_any).
func CanViewOrder(orderUserID, requester uuid.UUID, role string) bool {
	return IsOwnerOrHasPermission(orderUserID, requester, role, PermOrdersViewAny)
}

// CanUpdateOrderStatus encodes who may change order status and how.
// Roles with orders.manage_status may set any valid status. Customers may only cancel their own order from pending or approved.
func CanUpdateOrderStatus(orderUserID, requester uuid.UUID, role, currentStatus, newStatus string) bool {
	if HasPermission(role, PermOrdersManageStatus) {
		return true
	}
	if orderUserID != requester {
		return false
	}
	if newStatus != "cancelled" {
		return false
	}
	return currentStatus == "pending" || currentStatus == "approved"
}

// CanAccessConfiguration reports read/update access to a configuration record.
func CanAccessConfiguration(ownerUserID, requester uuid.UUID, role string) bool {
	return IsOwnerOrHasPermission(ownerUserID, requester, role, PermConfigurationsViewAny)
}

// CustomerMayChangeConfigurationStatus limits which status transitions users may trigger themselves.
func CustomerMayChangeConfigurationStatus(fromStatus, toStatus string) bool {
	switch toStatus {
	case "confirmed":
		return fromStatus == "draft"
	case "cancelled":
		return fromStatus == "draft" || fromStatus == "confirmed"
	default:
		return false
	}
}

// CanManageConfigurationStatus reports whether staff may set an arbitrary configuration status.
func CanManageConfigurationStatus(role string) bool {
	return HasPermission(role, PermConfigurationsManage)
}
