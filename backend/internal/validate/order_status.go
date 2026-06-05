package validate

import (
	"regexp"
	"strings"
)

const (
	OrderStatusCodeMax          = 32
	OrderStatusLabelMax         = 120
	OrderStatusDescriptionMax   = 2000
	OrderStatusSortMin          = -10_000
	OrderStatusSortMax          = 100_000
)

var orderStatusCodePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// OrderStatusCode validates machine-readable status code.
func OrderStatusCode(code string) (string, string) {
	s := strings.TrimSpace(strings.ToLower(code))
	if s == "" {
		return "", "code is required"
	}
	if runeLen(s) > OrderStatusCodeMax {
		return "", "code is too long (max 32 characters)"
	}
	if !orderStatusCodePattern.MatchString(s) {
		return "", "invalid order status code"
	}
	return s, ""
}

// OrderStatusCustomerLabel validates required customer-facing label.
func OrderStatusCustomerLabel(label string) (string, string) {
	return requiredSingleLine("customer_label_ru", label, OrderStatusLabelMax)
}

// OrderStatusAdminLabel validates optional admin label.
func OrderStatusAdminLabel(label *string) (*string, string) {
	return optionalSingleLine("admin_label_ru", label, OrderStatusLabelMax)
}

// OrderStatusDescription validates optional status description.
func OrderStatusDescription(description *string) (*string, string) {
	return optionalMultiline("description", description, OrderStatusDescriptionMax)
}

// OrderStatusSortOrder validates sort order range.
func OrderStatusSortOrder(sortOrder int) string {
	if sortOrder < OrderStatusSortMin || sortOrder > OrderStatusSortMax {
		return "invalid sort_order"
	}
	return ""
}
