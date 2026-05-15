package validate

import (
	"testing"
	"time"
)

func TestAppointmentDate_PastRejected(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	past := now.Add(-1 * time.Hour)
	if err := AppointmentDate(past, now); err == nil {
		t.Fatal("expected past appointment to fail")
	}
}

func TestAppointmentDate_FutureOK(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	future := now.Add(24 * time.Hour)
	if err := AppointmentDate(future, now); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestAppointmentDescription_TooLong(t *testing.T) {
	long := make([]rune, AppointmentDescriptionMaxRunes+1)
	for i := range long {
		long[i] = 'a'
	}
	s := string(long)
	if err := AppointmentDescription(&s); err == nil {
		t.Fatal("expected too long description error")
	}
}
