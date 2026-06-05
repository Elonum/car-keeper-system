package validate

import (
	"strings"
	"testing"
)

func TestBrandName_Valid(t *testing.T) {
	got, msg := BrandName("  Toyota ")
	if msg != "" || got != "Toyota" {
		t.Fatalf("got %q msg %q", got, msg)
	}
}

func TestBrandName_TooLong(t *testing.T) {
	_, msg := BrandName(strings.Repeat("a", BrandNameMax+1))
	if msg == "" {
		t.Fatal("expected error")
	}
}

func TestServiceCategory_Invalid(t *testing.T) {
	_, msg := ServiceCategory("unknown")
	if msg != "invalid service category" {
		t.Fatalf("got %q", msg)
	}
}

func TestServicePrice_Invalid(t *testing.T) {
	if msg := ServicePrice(-1); msg != "invalid price" {
		t.Fatalf("got %q", msg)
	}
}

func TestServiceDurationMinutes_OutOfRange(t *testing.T) {
	d := 0
	if msg := ServiceDurationMinutes(&d); msg == "" {
		t.Fatal("expected error")
	}
}
