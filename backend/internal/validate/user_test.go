package validate

import "testing"

func TestNames_Valid(t *testing.T) {
	fn, ln, msg := Names(" Иван ", "Петров")
	if msg != "" {
		t.Fatalf("unexpected: %s", msg)
	}
	if fn != "Иван" || ln != "Петров" {
		t.Fatalf("trim failed: %q %q", fn, ln)
	}
}

func TestNames_EmptyFirst(t *testing.T) {
	_, _, msg := Names("", "Петров")
	if msg == "" {
		t.Fatal("expected error")
	}
}

func TestNewPassword_TooShort(t *testing.T) {
	if msg := NewPassword("12345"); msg == "" {
		t.Fatal("expected min length error")
	}
}

func TestNewPassword_Valid(t *testing.T) {
	if msg := NewPassword("secret12"); msg != "" {
		t.Fatalf("unexpected: %s", msg)
	}
}

func TestPhonePtr_EmptyBecomesNil(t *testing.T) {
	empty := ""
	got, msg := PhonePtr(&empty)
	if msg != "" || got != nil {
		t.Fatalf("got %v msg %q", got, msg)
	}
}
