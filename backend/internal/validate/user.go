// Package validate holds shared string validation for user-related API input.
// Length limits match PostgreSQL varchar columns (rune count, not UTF-8 bytes).
package validate

import (
	"strings"
	"unicode/utf8"
)

const (
	MaxNameRunes     = 100
	MaxEmailRunes    = 255
	MaxPhoneRunes    = 30
	MinPasswordRunes = 6
	MaxPasswordRunes = 128
)

// Names trims and validates first/last name; returns error message or "".
func Names(firstName, lastName string) (fn, ln string, errMsg string) {
	fn = strings.TrimSpace(firstName)
	ln = strings.TrimSpace(lastName)
	if fn == "" || utf8.RuneCountInString(fn) > MaxNameRunes {
		return "", "", "first_name is required (max 100 characters)"
	}
	if ln == "" || utf8.RuneCountInString(ln) > MaxNameRunes {
		return "", "", "last_name is required (max 100 characters)"
	}
	return fn, ln, ""
}

// Email normalizes and validates email; returns error message or "".
func Email(email string) (normalized string, errMsg string) {
	s := strings.TrimSpace(strings.ToLower(email))
	if s == "" || utf8.RuneCountInString(s) > MaxEmailRunes {
		return "", "email is required (max 255 characters)"
	}
	return s, ""
}

// PhonePtr normalizes optional phone: empty string becomes nil; returns error message or "".
func PhonePtr(phone *string) (*string, string) {
	if phone == nil {
		return nil, ""
	}
	p := strings.TrimSpace(*phone)
	if p == "" {
		return nil, ""
	}
	if utf8.RuneCountInString(p) > MaxPhoneRunes {
		return nil, "phone is too long (max 30 characters)"
	}
	return &p, ""
}

// NewPassword validates a new password (register or change); returns error message or "".
func NewPassword(password string) string {
	n := utf8.RuneCountInString(password)
	if n < MinPasswordRunes {
		return "password must be at least 6 characters"
	}
	if n > MaxPasswordRunes {
		return "password is too long (max 128 characters)"
	}
	return ""
}

// NewPasswordMustDifferFromCurrent returns an error message if new equals the current password.
// Call only after the current password has been verified against the stored hash.
func NewPasswordMustDifferFromCurrent(currentPlain, newPlain string) string {
	if currentPlain == newPlain {
		return "new password must be different from the current password"
	}
	return ""
}
