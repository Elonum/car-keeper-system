package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DocumentRepository struct {
	db *database.DB
}

func NewDocumentRepository(db *database.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Insert(ctx context.Context, d model.Document) error {
	query := `
		INSERT INTO documents (
			document_id, user_id, order_id, service_appointment_id,
			document_type, file_path, file_name, file_size, mime_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		d.DocumentID, d.UserID, d.OrderID, d.ServiceAppointmentID,
		d.DocumentType, d.FilePath, d.FileName, d.FileSize, d.MimeType,
	)
	if err != nil {
		return fmt.Errorf("insert document: %w", err)
	}
	return nil
}

func (r *DocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Document, error) {
	var d model.Document
	query := `
		SELECT document_id, user_id, order_id, service_appointment_id,
		       document_type, file_path, file_name, file_size, mime_type, created_at
		FROM documents WHERE document_id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&d.DocumentID, &d.UserID, &d.OrderID, &d.ServiceAppointmentID,
		&d.DocumentType, &d.FilePath, &d.FileName, &d.FileSize, &d.MimeType, &d.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("get document: %w", err)
	}
	return &d, nil
}

func (r *DocumentRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Document, error) {
	query := `
		SELECT document_id, user_id, order_id, service_appointment_id,
		       document_type, file_path, file_name, file_size, mime_type, created_at
		FROM documents WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()
	return scanDocuments(rows)
}

func (r *DocumentRepository) ListByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Document, error) {
	query := `
		SELECT document_id, user_id, order_id, service_appointment_id,
		       document_type, file_path, file_name, file_size, mime_type, created_at
		FROM documents WHERE order_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()
	return scanDocuments(rows)
}

func (r *DocumentRepository) ListByServiceAppointmentID(ctx context.Context, apptID uuid.UUID) ([]model.Document, error) {
	query := `
		SELECT document_id, user_id, order_id, service_appointment_id,
		       document_type, file_path, file_name, file_size, mime_type, created_at
		FROM documents WHERE service_appointment_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, apptID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()
	return scanDocuments(rows)
}

func scanDocuments(rows pgx.Rows) ([]model.Document, error) {
	var list []model.Document
	for rows.Next() {
		var d model.Document
		if err := rows.Scan(
			&d.DocumentID, &d.UserID, &d.OrderID, &d.ServiceAppointmentID,
			&d.DocumentType, &d.FilePath, &d.FileName, &d.FileSize, &d.MimeType, &d.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

func (r *DocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.db.Pool.Exec(ctx, `DELETE FROM documents WHERE document_id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("document not found")
	}
	return nil
}
