package authz

import (
	"github.com/google/uuid"
)

// IsStaff returns true for roles that may manage orders, news, and see unpublished content.
func IsStaff(role string) bool {
	return role == "admin" || role == "manager"
}

// IsOwnerOrStaff reports whether the requester owns the resource or has a staff role.
func IsOwnerOrStaff(resourceOwnerID, requester uuid.UUID, role string) bool {
	if IsStaff(role) {
		return true
	}
	return resourceOwnerID == requester
}

// CanViewOrder reports whether the requester may read the order (owner or staff).
func CanViewOrder(orderUserID, requester uuid.UUID, role string) bool {
	return IsOwnerOrStaff(orderUserID, requester, role)
}

// CanUpdateOrderStatus encodes who may change order status and how.
// Staff may set any valid status. Customers may only cancel their own order from pending or approved.
func CanUpdateOrderStatus(orderUserID, requester uuid.UUID, role, currentStatus, newStatus string) bool {
	if IsStaff(role) {
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
	return IsOwnerOrStaff(ownerUserID, requester, role)
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
