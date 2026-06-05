package validate

import (
	"strings"
	"testing"
)

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

func TestNames_InvalidCharacters(t *testing.T) {
	_, _, msg := Names("John<script>", "Doe")
	if msg != "first_name contains invalid characters" {
		t.Fatalf("got %q", msg)
	}
	_, _, msg = Names("John", "Doe123")
	if msg != "last_name contains invalid characters" {
		t.Fatalf("got %q", msg)
	}
}

func TestEmail_Valid(t *testing.T) {
	em, msg := Email("  User@Example.COM ")
	if msg != "" {
		t.Fatalf("unexpected: %s", msg)
	}
	if em != "user@example.com" {
		t.Fatalf("normalize failed: %q", em)
	}
}

func TestEmail_InvalidFormat(t *testing.T) {
	_, msg := Email("not-an-email")
	if msg != "invalid email format" {
		t.Fatalf("got %q", msg)
	}
}

func TestEmail_TooLong(t *testing.T) {
	longLocal := strings.Repeat("a", MaxEmailRunes)
	_, msg := Email(longLocal + "@x.com")
	if msg == "" {
		t.Fatal("expected max length error")
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

func TestPhonePtr_ValidNormalized(t *testing.T) {
	raw := "+7 (999) 123-45-67"
	got, msg := PhonePtr(&raw)
	if msg != "" {
		t.Fatalf("unexpected: %s", msg)
	}
	if got == nil || *got != "+79991234567" {
		t.Fatalf("got %v", got)
	}
}

func TestPhonePtr_InvalidCharacters(t *testing.T) {
	raw := "abc1234567890"
	got, msg := PhonePtr(&raw)
	if msg != "phone contains invalid characters" || got != nil {
		t.Fatalf("got %v msg %q", got, msg)
	}
}

func TestPhonePtr_TooFewDigits(t *testing.T) {
	raw := "+7 999"
	got, msg := PhonePtr(&raw)
	if msg != "phone must contain 10–15 digits" || got != nil {
		t.Fatalf("got %v msg %q", got, msg)
	}
}

func TestPhonePtr_TooManyDigits(t *testing.T) {
	raw := strings.Repeat("1", 16)
	got, msg := PhonePtr(&raw)
	if msg != "phone must contain 10–15 digits" || got != nil {
		t.Fatalf("got %v msg %q", got, msg)
	}
}

func TestPhonePtr_TooLongFormatted(t *testing.T) {
	raw := "+1 " + strings.Repeat("(9) ", 20) + strings.Repeat("9", 10)
	got, msg := PhonePtr(&raw)
	if msg != "phone is too long (max 30 characters)" || got != nil {
		t.Fatalf("got %v msg %q", got, msg)
	}
}
