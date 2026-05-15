package service

import (
	"testing"
	"time"

	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
)

func TestTotalDurationMinutes_DefaultWhenEmpty(t *testing.T) {
	if got := totalDurationMinutes(nil); got != 30 {
		t.Fatalf("got %d", got)
	}
}

func TestTotalDurationMinutes_Sum(t *testing.T) {
	d30 := 30
	d60 := 60
	types := []model.ServiceType{
		{DurationMinutes: &d30},
		{DurationMinutes: &d60},
	}
	if got := totalDurationMinutes(types); got != 90 {
		t.Fatalf("got %d", got)
	}
}

func TestDedupeUUIDs(t *testing.T) {
	id := uuid.New()
	got := dedupeUUIDs([]uuid.UUID{id, id, uuid.New()})
	if len(got) != 2 {
		t.Fatalf("got %d", len(got))
	}
}

func TestValidateAppointmentSlot_OutsideWorkday(t *testing.T) {
	branch := &model.Branch{
		Timezone:            "Europe/Moscow",
		WorkdayStartMinutes: 540,
		WorkdayEndMinutes:   1080,
		SlotStepMinutes:     30,
		ConcurrentBays:      2,
	}
	loc, _ := time.LoadLocation("Europe/Moscow")
	// Monday 2026-05-18 06:00 local — before workday
	start := time.Date(2026, 5, 18, 6, 0, 0, 0, loc)
	if err := validateAppointmentSlot(branch, start, 30); err == nil {
		t.Fatal("expected outside working hours error")
	}
}
