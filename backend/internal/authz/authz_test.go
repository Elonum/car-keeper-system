package authz

import (
	"testing"

	"github.com/google/uuid"
)

func TestIsStaff(t *testing.T) {
	if !IsStaff("admin") || !IsStaff("manager") {
		t.Fatal("expected admin and manager to be staff")
	}
	if IsStaff("customer") || IsStaff("") {
		t.Fatal("expected non-staff roles")
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
		t.Fatal("staff may set any transition")
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

func TestIsOwnerOrStaff(t *testing.T) {
	a, b := uuid.New(), uuid.New()
	if !IsOwnerOrStaff(a, a, "customer") {
		t.Fatal("owner")
	}
	if IsOwnerOrStaff(a, b, "customer") {
		t.Fatal("stranger")
	}
	if !IsOwnerOrStaff(a, b, "admin") {
		t.Fatal("staff")
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
