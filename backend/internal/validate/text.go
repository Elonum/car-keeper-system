package validate

import (
	"strings"
	"unicode/utf8"
)

func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}

func hasDisallowedControlRunes(s string, allowNewlines bool) bool {
	for _, r := range s {
		if allowNewlines && (r == '\n' || r == '\r' || r == '\t') {
			continue
		}
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}

func requiredSingleLine(field, value string, maxRunes int) (string, string) {
	s := strings.TrimSpace(value)
	if s == "" {
		return "", field + " is required"
	}
	if runeLen(s) > maxRunes {
		return "", field + " is too long (max " + itoa(maxRunes) + " characters)"
	}
	if hasDisallowedControlRunes(s, false) {
		return "", field + " contains invalid characters"
	}
	return s, ""
}

func optionalSingleLine(field string, value *string, maxRunes int) (*string, string) {
	if value == nil {
		return nil, ""
	}
	s := strings.TrimSpace(*value)
	if s == "" {
		return nil, ""
	}
	if runeLen(s) > maxRunes {
		return nil, field + " is too long (max " + itoa(maxRunes) + " characters)"
	}
	if hasDisallowedControlRunes(s, false) {
		return nil, field + " contains invalid characters"
	}
	return &s, ""
}

func optionalMultiline(field string, value *string, maxRunes int) (*string, string) {
	if value == nil {
		return nil, ""
	}
	s := strings.TrimSpace(*value)
	if s == "" {
		return nil, ""
	}
	if runeLen(s) > maxRunes {
		return nil, field + " is too long (max " + itoa(maxRunes) + " characters)"
	}
	if hasDisallowedControlRunes(s, true) {
		return nil, field + " contains invalid characters"
	}
	return &s, ""
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
