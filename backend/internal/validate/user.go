// Package validate holds shared string validation for user-related API input.
// Length limits match PostgreSQL varchar columns (rune count, not UTF-8 bytes).
package validate

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	MaxNameRunes     = 100
	MaxEmailRunes    = 255
	MaxPhoneRunes    = 30
	MinPhoneDigits   = 10
	MaxPhoneDigits   = 15
	MinPasswordRunes = 6
	MaxPasswordRunes = 128
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func isValidNameRune(r rune) bool {
	return unicode.IsLetter(r) || r == ' ' || r == '-' || r == '\'' || r == '.'
}

func nameCharsOK(s string) bool {
	for _, r := range s {
		if !isValidNameRune(r) {
			return false
		}
	}
	return true
}

func phoneDigitCount(s string) int {
	n := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			n++
		}
	}
	return n
}

func phoneCharsOK(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			continue
		}
		switch r {
		case '+', ' ', '(', ')', '-', '.':
			continue
		default:
			return false
		}
	}
	return true
}

func normalizePhone(s string) string {
	trimmed := strings.TrimSpace(s)
	var digits strings.Builder
	hasPlus := false
	for _, r := range trimmed {
		if r == '+' {
			hasPlus = true
			continue
		}
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	d := digits.String()
	if d == "" {
		return trimmed
	}
	if hasPlus {
		return "+" + d
	}
	return d
}

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
	if !nameCharsOK(fn) {
		return "", "", "first_name contains invalid characters"
	}
	if !nameCharsOK(ln) {
		return "", "", "last_name contains invalid characters"
	}
	return fn, ln, ""
}

// Email normalizes and validates email; returns error message or "".
func Email(email string) (normalized string, errMsg string) {
	s := strings.TrimSpace(strings.ToLower(email))
	if s == "" || utf8.RuneCountInString(s) > MaxEmailRunes {
		return "", "email is required (max 255 characters)"
	}
	if !emailPattern.MatchString(s) {
		return "", "invalid email format"
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
	if !phoneCharsOK(p) {
		return nil, "phone contains invalid characters"
	}
	digits := phoneDigitCount(p)
	if digits < MinPhoneDigits || digits > MaxPhoneDigits {
		return nil, "phone must contain 10–15 digits"
	}
	norm := normalizePhone(p)
	if utf8.RuneCountInString(norm) > MaxPhoneRunes {
		return nil, "phone is too long (max 30 characters)"
	}
	return &norm, ""
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
