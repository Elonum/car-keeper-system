package validate

import (
	"math"
	"strings"
)

const (
	BrandNameMax          = 150
	BrandCountryMax       = 100
	ModelNameMax          = 150
	ModelSegmentMax       = 100
	ModelDescriptionMax   = 2000
	ServiceTypeNameMax    = 150
	ServiceDescriptionMax = 2000
	BranchNameMax         = 200
	BranchAddressMax      = 2000
	ServicePriceMax       = 99_999_999.99
	ServiceDurationMin    = 1
	ServiceDurationMax    = 1440
)

var serviceCategories = map[string]struct{}{
	"maintenance":   {},
	"repair":        {},
	"diagnostics":   {},
	"detailing":     {},
	"tires":         {},
}

// BrandName validates a car brand name.
func BrandName(name string) (string, string) {
	return requiredSingleLine("name", name, BrandNameMax)
}

// BrandCountry validates brand country of origin.
func BrandCountry(country string) (string, string) {
	return requiredSingleLine("country", country, BrandCountryMax)
}

// ModelName validates a car model name.
func ModelName(name string) (string, string) {
	return requiredSingleLine("name", name, ModelNameMax)
}

// ModelSegment validates optional model segment.
func ModelSegment(segment *string) (*string, string) {
	return optionalSingleLine("segment", segment, ModelSegmentMax)
}

// ModelDescription validates optional model description.
func ModelDescription(description *string) (*string, string) {
	return optionalMultiline("description", description, ModelDescriptionMax)
}

// ServiceTypeName validates a service offering name.
func ServiceTypeName(name string) (string, string) {
	return requiredSingleLine("name", name, ServiceTypeNameMax)
}

// ServiceCategory validates service category against DB constraint.
func ServiceCategory(category string) (string, string) {
	s, msg := requiredSingleLine("category", category, 50)
	if msg != "" {
		return "", msg
	}
	if _, ok := serviceCategories[s]; !ok {
		return "", "invalid service category"
	}
	return s, ""
}

// ServiceDescription validates optional service description.
func ServiceDescription(description *string) (*string, string) {
	return optionalMultiline("description", description, ServiceDescriptionMax)
}

// ServicePrice validates non-negative catalog price within numeric(12,2).
func ServicePrice(price float64) string {
	if math.IsNaN(price) || math.IsInf(price, 0) {
		return "invalid price"
	}
	if price < 0 {
		return "invalid price"
	}
	if price > ServicePriceMax {
		return "price is too large"
	}
	return ""
}

// ServiceDurationMinutes validates optional service duration.
func ServiceDurationMinutes(duration *int) string {
	if duration == nil {
		return ""
	}
	if *duration < ServiceDurationMin || *duration > ServiceDurationMax {
		return "duration_minutes must be between 1 and 1440"
	}
	return ""
}

// BranchName validates branch display name on update.
func BranchName(name string) (string, string) {
	return requiredSingleLine("name", name, BranchNameMax)
}

// BranchAddress validates branch address on update.
func BranchAddress(address string) (string, string) {
	s := strings.TrimSpace(address)
	if s == "" {
		return "", "address is required"
	}
	if runeLen(s) > BranchAddressMax {
		return "", "address is too long (max 2000 characters)"
	}
	if hasDisallowedControlRunes(s, true) {
		return "", "address contains invalid characters"
	}
	return s, ""
}

// EmailOptionalPtr normalizes optional email; empty becomes nil.
func EmailOptionalPtr(email *string) (*string, string) {
	if email == nil {
		return nil, ""
	}
	s := strings.TrimSpace(*email)
	if s == "" {
		return nil, ""
	}
	normalized, msg := Email(s)
	if msg != "" {
		return nil, msg
	}
	return &normalized, ""
}
