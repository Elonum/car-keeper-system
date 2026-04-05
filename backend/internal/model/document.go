package model

import (
	"time"

	"github.com/google/uuid"
)

// Document matches table documents (file_path is internal, not exposed in JSON).
type Document struct {
	DocumentID           uuid.UUID  `db:"document_id" json:"document_id"`
	UserID               uuid.UUID  `db:"user_id" json:"user_id"`
	OrderID              *uuid.UUID `db:"order_id" json:"order_id,omitempty"`
	ServiceAppointmentID *uuid.UUID `db:"service_appointment_id" json:"service_appointment_id,omitempty"`
	DocumentType         string     `db:"document_type" json:"document_type"`
	FilePath             string     `db:"file_path" json:"-"`
	FileName             *string    `db:"file_name" json:"file_name,omitempty"`
	FileSize             *int64     `db:"file_size" json:"file_size,omitempty"`
	MimeType             *string    `db:"mime_type" json:"mime_type,omitempty"`
	CreatedAt            time.Time  `db:"created_at" json:"created_at"`
}

// DocumentTypes lists allowed document_type values (must match DB CHECK).
var DocumentTypes = map[string]struct{}{
	"commercial_offer": {},
	"order_contract":   {},
	"service_order":    {},
	"service_act":      {},
}

func ValidDocumentType(t string) bool {
	_, ok := DocumentTypes[t]
	return ok
}
