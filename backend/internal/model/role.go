package model

import (
	"time"

	"github.com/google/uuid"
)

// RoleDefinition is a row in role_definitions (stable code for JWT and policies).
type RoleDefinition struct {
	RoleID      uuid.UUID `db:"role_id" json:"role_id"`
	Code        string    `db:"code" json:"code"`
	NameRu      string    `db:"name_ru" json:"name_ru"`
	Description *string   `db:"description" json:"description,omitempty"`
	SortOrder   int       `db:"sort_order" json:"sort_order"`
	IsStaff     bool      `db:"is_staff" json:"is_staff"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
