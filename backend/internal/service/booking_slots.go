package service

import (
	"context"
	"fmt"
	"time"

	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
)

func totalDurationMinutes(types []model.ServiceType) int {
	sum := 0
	for _, t := range types {
		d := 30
		if t.DurationMinutes != nil && *t.DurationMinutes > 0 {
			d = *t.DurationMinutes
		}
		sum += d
	}
	if sum < 1 {
		return 30
	}
	return sum
}

func validateAppointmentSlot(branch *model.Branch, start time.Time, durationMin int) error {
	if durationMin <= 0 {
		return fmt.Errorf("invalid service duration")
	}
	loc, err := time.LoadLocation(branch.Timezone)
	if err != nil {
		loc = time.UTC
	}
	local := start.In(loc)
	if local.Weekday() == time.Sunday {
		return fmt.Errorf("branch is closed on this day")
	}
	dayStart := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	minsFromMidnight := int(local.Sub(dayStart).Minutes())
	if minsFromMidnight < branch.WorkdayStartMinutes || minsFromMidnight+durationMin > branch.WorkdayEndMinutes {
		return fmt.Errorf("appointment outside branch working hours")
	}
	if (minsFromMidnight-branch.WorkdayStartMinutes)%branch.SlotStepMinutes != 0 {
		return fmt.Errorf("invalid appointment time slot")
	}
	return nil
}

// BranchAvailability returns concrete UTC start times for which booking is possible on a calendar day.
func (s *ServiceService) BranchAvailability(ctx context.Context, branchID uuid.UUID, dateStr string, serviceTypeIDs []uuid.UUID) (*model.BranchAvailability, error) {
	branch, err := s.repo.Branch.GetByID(ctx, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}
	if !branch.IsActive {
		return nil, fmt.Errorf("branch is not active")
	}

	uniq := dedupeUUIDs(serviceTypeIDs)
	if len(uniq) == 0 {
		return nil, fmt.Errorf("at least one service type is required")
	}

	types, err := s.repo.ServiceType.GetByIDs(ctx, uniq)
	if err != nil {
		return nil, err
	}
	if len(types) != len(uniq) {
		return nil, fmt.Errorf("one or more service types are not available")
	}
	duration := totalDurationMinutes(types)

	loc, err := time.LoadLocation(branch.Timezone)
	if err != nil {
		loc, _ = time.LoadLocation("Europe/Moscow")
	}
	dayStart, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}
	if dayStart.Weekday() == time.Sunday {
		return &model.BranchAvailability{
			SlotStarts:      nil,
			Timezone:        branch.Timezone,
			DurationMinutes: duration,
		}, nil
	}

	workStart := dayStart.Add(time.Duration(branch.WorkdayStartMinutes) * time.Minute)
	workEnd := dayStart.Add(time.Duration(branch.WorkdayEndMinutes) * time.Minute)
	window := int(workEnd.Sub(workStart).Minutes())
	if duration > window {
		return &model.BranchAvailability{
			SlotStarts:      nil,
			Timezone:        branch.Timezone,
			DurationMinutes: duration,
		}, nil
	}

	step := time.Duration(branch.SlotStepMinutes) * time.Minute
	now := time.Now()
	minLead := 15 * time.Minute

	var slots []time.Time
	for t := workStart; ; t = t.Add(step) {
		tEnd := t.Add(time.Duration(duration) * time.Minute)
		if tEnd.After(workEnd) {
			break
		}
		if sameDay(t, now.In(loc)) && t.Before(now.In(loc).Add(minLead)) {
			continue
		}
		n, err := s.repo.ServiceAppointment.CountOverlappingScheduled(ctx, branchID, t, tEnd)
		if err != nil {
			return nil, err
		}
		if n < branch.ConcurrentBays {
			slots = append(slots, t.UTC())
		}
	}

	return &model.BranchAvailability{
		SlotStarts:      slots,
		Timezone:        branch.Timezone,
		DurationMinutes: duration,
	}, nil
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func dedupeUUIDs(ids []uuid.UUID) []uuid.UUID {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[uuid.UUID]struct{}, len(ids))
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
