package authz

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	SetRolePermissions(DefaultRolePermissions())
	SetStaffRoles([]string{"admin", "manager", "service_advisor"})
	os.Exit(m.Run())
}

func TestIsStaff(t *testing.T) {
	if !IsStaff("admin") || !IsStaff("manager") || !IsStaff("service_advisor") {
		t.Fatal("expected admin, manager, service_advisor to be staff")
	}
	if IsStaff("customer") || IsStaff("") {
		t.Fatal("expected non-staff roles")
	}
}

func TestHasPermissionMatrix(t *testing.T) {
	if !HasPermission("manager", PermNewsManage) {
		t.Fatal("manager should manage news")
	}
	if HasPermission("service_advisor", PermNewsManage) {
		t.Fatal("service advisor should not manage news")
	}
	if !HasPermission("service_advisor", PermOrdersManageStatus) {
		t.Fatal("service advisor should manage order status")
	}
	if !HasPermission("admin", PermAdminOrderStatuses) {
		t.Fatal("admin should CRUD order statuses")
	}
	if HasPermission("manager", PermAdminOrderStatuses) {
		t.Fatal("manager should not CRUD order status definitions")
	}
	if !HasPermission("admin", PermCatalogManage) || !HasPermission("admin", PermServiceManage) {
		t.Fatal("admin should have catalog and service manage")
	}
	if HasPermission("manager", PermCatalogManage) {
		t.Fatal("manager should not manage full catalog")
	}
	if !HasPermission("manager", PermServiceManage) {
		t.Fatal("manager should manage service catalog")
	}
}

func TestIsAdmin(t *testing.T) {
	if !IsAdmin("admin") {
		t.Fatal("admin")
	}
	if IsAdmin("manager") || IsAdmin("service_advisor") || IsAdmin("customer") {
		t.Fatal("only admin")
	}
}

func TestCanViewOrder(t *testing.T) {
	owner := uuid.New()
	other := uuid.New()
	if !CanViewOrder(owner, owner, "customer") {
		t.Fatal("owner should view own order")
	}
	if CanViewOrder(owner, other, "customer") {
		t.Fatal("other customer should not view")
	}
	if !CanViewOrder(owner, other, "admin") {
		t.Fatal("admin should view any order")
	}
}

func TestCanUpdateOrderStatus(t *testing.T) {
	owner := uuid.New()
	other := uuid.New()
	if !CanUpdateOrderStatus(owner, owner, "customer", "pending", "cancelled") {
		t.Fatal("customer may cancel pending")
	}
	if CanUpdateOrderStatus(owner, owner, "customer", "paid", "cancelled") {
		t.Fatal("customer should not cancel paid order")
	}
	if CanUpdateOrderStatus(owner, owner, "customer", "pending", "paid") {
		t.Fatal("customer should not set paid")
	}
	if !CanUpdateOrderStatus(owner, other, "admin", "pending", "paid") {
		t.Fatal("staff with orders.manage_status may set any transition")
	}
}

func TestCanAccessConfiguration(t *testing.T) {
	u1, u2 := uuid.New(), uuid.New()
	if !CanAccessConfiguration(u1, u1, "customer") {
		t.Fatal("owner access")
	}
	if CanAccessConfiguration(u1, u2, "customer") {
		t.Fatal("no foreign access")
	}
	if !CanAccessConfiguration(u1, u2, "manager") {
		t.Fatal("staff access")
	}
}

func TestIsOwnerOrHasPermission(t *testing.T) {
	a, b := uuid.New(), uuid.New()
	if !IsOwnerOrHasPermission(a, a, "customer", PermAppointmentsViewAny) {
		t.Fatal("owner")
	}
	if IsOwnerOrHasPermission(a, b, "customer", PermAppointmentsViewAny) {
		t.Fatal("stranger")
	}
	if !IsOwnerOrHasPermission(a, b, "admin", PermAppointmentsViewAny) {
		t.Fatal("staff with permission")
	}
}

func TestCustomerMayChangeConfigurationStatus(t *testing.T) {
	if !CustomerMayChangeConfigurationStatus("draft", "confirmed") {
		t.Fatal("draft -> confirmed")
	}
	if CustomerMayChangeConfigurationStatus("draft", "ordered") {
		t.Fatal("draft -> ordered blocked for customer")
	}
	if !CustomerMayChangeConfigurationStatus("confirmed", "cancelled") {
		t.Fatal("confirmed -> cancelled")
	}
}
