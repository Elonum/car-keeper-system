package model

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	BranchID            uuid.UUID `db:"branch_id" json:"branch_id"`
	Name                string    `db:"name" json:"name"`
	Address             string    `db:"address" json:"address"`
	Phone               *string   `db:"phone" json:"phone,omitempty"`
	Email               *string   `db:"email" json:"email,omitempty"`
	IsActive            bool      `db:"is_active" json:"is_active"`
	Timezone            string    `db:"timezone" json:"timezone"`
	WorkdayStartMinutes int       `db:"workday_start_minutes" json:"workday_start_minutes"`
	WorkdayEndMinutes   int       `db:"workday_end_minutes" json:"workday_end_minutes"`
	SlotStepMinutes     int       `db:"slot_step_minutes" json:"slot_step_minutes"`
	ConcurrentBays      int       `db:"concurrent_bays" json:"concurrent_bays"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

