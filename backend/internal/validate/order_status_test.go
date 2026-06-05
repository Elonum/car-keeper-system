package validate

import "testing"

func TestOrderStatusCode_Valid(t *testing.T) {
	got, msg := OrderStatusCode(" Pending_Payment ")
	if msg != "" || got != "pending_payment" {
		t.Fatalf("got %q msg %q", got, msg)
	}
}

func TestOrderStatusCode_Invalid(t *testing.T) {
	_, msg := OrderStatusCode("9bad")
	if msg != "invalid order status code" {
		t.Fatalf("got %q", msg)
	}
}

func TestOrderStatusSortOrder_Invalid(t *testing.T) {
	if msg := OrderStatusSortOrder(OrderStatusSortMax + 1); msg == "" {
		t.Fatal("expected error")
	}
}
