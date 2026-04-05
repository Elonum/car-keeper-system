package validate

import (
	"fmt"
	"time"
	"unicode/utf8"
)

// AppointmentDescriptionMaxRunes aligns with frontend UI_LIMITS.APPOINTMENT_DESCRIPTION.
const AppointmentDescriptionMaxRunes = 10000

// AppointmentDate checks that the visit is not in the past (small clock skew) and not unreasonably far ahead.
func AppointmentDate(t time.Time, now time.Time) error {
	if t.IsZero() {
		return fmt.Errorf("appointment_date is required")
	}
	tu := t.UTC()
	nu := now.UTC()
	if !tu.After(nu.Add(-2 * time.Minute)) {
		return fmt.Errorf("appointment must be in the future")
	}
	max := nu.AddDate(1, 0, 1)
	if tu.After(max) {
		return fmt.Errorf("appointment date is too far in the future")
	}
	return nil
}

// AppointmentDescription limits payload size for TEXT column abuse.
func AppointmentDescription(p *string) error {
	if p == nil {
		return nil
	}
	if utf8.RuneCountInString(*p) > AppointmentDescriptionMaxRunes {
		return fmt.Errorf("description is too long")
	}
	return nil
}
