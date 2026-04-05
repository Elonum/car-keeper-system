package validate

import (
	"regexp"
	"strings"
)

// VIN pattern: 17 chars, no I, O, Q (ISO 3779).
var vinPattern = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)

// NormalizeVIN trims and uppercases VIN input.
func NormalizeVIN(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

// VIN returns an API error message if invalid, or "" if OK (after normalization).
func VIN(s string) string {
	v := NormalizeVIN(s)
	if len(v) != 17 {
		return "VIN must be exactly 17 characters"
	}
	if !vinPattern.MatchString(v) {
		return "VIN must contain only letters A–Z (except I, O, Q) and digits"
	}
	return ""
}
