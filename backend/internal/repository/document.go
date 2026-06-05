package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DocumentRepository struct {
	db *database.DB
}

const documentSelectWithContext = `
	SELECT
		d.document_id,
		d.user_id,
		d.order_id,
		d.service_appointment_id,
		d.document_type,
		d.file_path,
		d.file_name,
		d.file_size,
		d.mime_type,
		d.created_at,
		NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), '') AS owner_name,
		u.email::text AS owner_email,
		CASE
			WHEN d.order_id IS NOT NULL THEN 'order'
			WHEN d.service_appointment_id IS NOT NULL THEN 'service_appointment'
		END AS attachment_kind,
		CASE
			WHEN d.order_id IS NOT NULL THEN
				CONCAT(
					'Заказ #', LEFT(d.order_id::text, 8),
					COALESCE(' · ' || osd.customer_label_ru, '')
				)
			WHEN d.service_appointment_id IS NOT NULL THEN
				CONCAT(
					'Запись на ТО #', LEFT(d.service_appointment_id::text, 8),
					COALESCE(' · ' || TO_CHAR(sa.appointment_date AT TIME ZONE 'Europe/Moscow', 'DD.MM.YYYY HH24:MI'), ''),
					COALESCE(' · ' || b.name, ''),
					COALESCE(' · VIN ' || uc.vin, '')
				)
		END AS attachment_label
	FROM documents d
	JOIN users u ON u.user_id = d.user_id
	LEFT JOIN orders o ON o.order_id = d.order_id
	LEFT JOIN order_status_definitions osd ON osd.code = o.status
	LEFT JOIN service_appointments sa ON sa.service_appointment_id = d.service_appointment_id
	LEFT JOIN branches b ON b.branch_id = sa.branch_id
	LEFT JOIN user_cars uc ON uc.user_car_id = sa.user_car_id
`

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
	query := documentSelectWithContext + ` WHERE d.document_id = $1`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&d.DocumentID, &d.UserID, &d.OrderID, &d.ServiceAppointmentID,
		&d.DocumentType, &d.FilePath, &d.FileName, &d.FileSize, &d.MimeType, &d.CreatedAt,
		&d.OwnerName, &d.OwnerEmail, &d.AttachmentKind, &d.AttachmentLabel,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w", apperr.ErrNotFound)
		}
		return nil, fmt.Errorf("get document: %w", err)
	}
	return &d, nil
}

func (r *DocumentRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Document, error) {
	query := documentSelectWithContext + ` WHERE d.user_id = $1 ORDER BY d.created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()
	return scanDocuments(rows)
}

func (r *DocumentRepository) ListAll(ctx context.Context) ([]model.Document, error) {
	query := documentSelectWithContext + ` ORDER BY d.created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()
	return scanDocuments(rows)
}

func (r *DocumentRepository) ListByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Document, error) {
	query := documentSelectWithContext + ` WHERE d.order_id = $1 ORDER BY d.created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()
	return scanDocuments(rows)
}

func (r *DocumentRepository) ListByServiceAppointmentID(ctx context.Context, apptID uuid.UUID) ([]model.Document, error) {
	query := documentSelectWithContext + ` WHERE d.service_appointment_id = $1 ORDER BY d.created_at DESC`
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
			&d.OwnerName, &d.OwnerEmail, &d.AttachmentKind, &d.AttachmentLabel,
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
		return fmt.Errorf("%w", apperr.ErrNotFound)
	}
	return nil
}
