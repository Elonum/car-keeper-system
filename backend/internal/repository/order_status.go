package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderStatusRepository struct {
	db *database.DB
}

func NewOrderStatusRepository(db *database.DB) *OrderStatusRepository {
	return &OrderStatusRepository{db: db}
}

func (r *OrderStatusRepository) ListActive(ctx context.Context) ([]model.OrderStatusDefinition, error) {
	query := `
		SELECT order_status_id, code, customer_label_ru, admin_label_ru, description,
		       sort_order, is_active, is_terminal, created_at, updated_at
		FROM order_status_definitions
		WHERE is_active = true
		ORDER BY sort_order, code
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	defer rows.Close()
	return scanOrderStatusRows(rows)
}

func (r *OrderStatusRepository) ListAll(ctx context.Context) ([]model.OrderStatusDefinition, error) {
	query := `
		SELECT order_status_id, code, customer_label_ru, admin_label_ru, description,
		       sort_order, is_active, is_terminal, created_at, updated_at
		FROM order_status_definitions
		ORDER BY sort_order, code
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	defer rows.Close()
	return scanOrderStatusRows(rows)
}

func scanOrderStatusRows(rows pgx.Rows) ([]model.OrderStatusDefinition, error) {
	var out []model.OrderStatusDefinition
	for rows.Next() {
		var d model.OrderStatusDefinition
		if err := rows.Scan(
			&d.OrderStatusID, &d.Code, &d.CustomerLabelRu, &d.AdminLabelRu, &d.Description,
			&d.SortOrder, &d.IsActive, &d.IsTerminal, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, apperr.Internal(err)
		}
		out = append(out, d)
	}
	return out, nil
}

func (r *OrderStatusRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.OrderStatusDefinition, error) {
	query := `
		SELECT order_status_id, code, customer_label_ru, admin_label_ru, description,
		       sort_order, is_active, is_terminal, created_at, updated_at
		FROM order_status_definitions
		WHERE order_status_id = $1
	`
	var d model.OrderStatusDefinition
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&d.OrderStatusID, &d.Code, &d.CustomerLabelRu, &d.AdminLabelRu, &d.Description,
		&d.SortOrder, &d.IsActive, &d.IsTerminal, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w", apperr.ErrNotFound)
		}
		return nil, apperr.Internal(err)
	}
	return &d, nil
}

func (r *OrderStatusRepository) GetByCode(ctx context.Context, code string) (*model.OrderStatusDefinition, error) {
	code = strings.TrimSpace(code)
	query := `
		SELECT order_status_id, code, customer_label_ru, admin_label_ru, description,
		       sort_order, is_active, is_terminal, created_at, updated_at
		FROM order_status_definitions
		WHERE code = $1
	`
	var d model.OrderStatusDefinition
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(
		&d.OrderStatusID, &d.Code, &d.CustomerLabelRu, &d.AdminLabelRu, &d.Description,
		&d.SortOrder, &d.IsActive, &d.IsTerminal, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w", apperr.ErrNotFound)
		}
		return nil, apperr.Internal(err)
	}
	return &d, nil
}

func (r *OrderStatusRepository) Create(ctx context.Context, in model.OrderStatusCreate) (*model.OrderStatusDefinition, error) {
	code := strings.TrimSpace(in.Code)
	if code == "" {
		return nil, apperr.BadRequest("code is required")
	}
	customer := strings.TrimSpace(in.CustomerLabelRu)
	if customer == "" {
		return nil, apperr.BadRequest("customer_label_ru is required")
	}
	sortOrder := 0
	if in.SortOrder != nil {
		sortOrder = *in.SortOrder
	}
	isActive := true
	if in.IsActive != nil {
		isActive = *in.IsActive
	}
	isTerminal := false
	if in.IsTerminal != nil {
		isTerminal = *in.IsTerminal
	}

	query := `
		INSERT INTO order_status_definitions
			(code, customer_label_ru, admin_label_ru, description, sort_order, is_active, is_terminal)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING order_status_id, code, customer_label_ru, admin_label_ru, description,
		          sort_order, is_active, is_terminal, created_at, updated_at
	`
	var d model.OrderStatusDefinition
	err := r.db.Pool.QueryRow(ctx, query,
		code, customer, in.AdminLabelRu, in.Description, sortOrder, isActive, isTerminal,
	).Scan(
		&d.OrderStatusID, &d.Code, &d.CustomerLabelRu, &d.AdminLabelRu, &d.Description,
		&d.SortOrder, &d.IsActive, &d.IsTerminal, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperr.Conflict("Order status code already exists")
		}
		return nil, apperr.Internal(err)
	}
	return &d, nil
}

func (r *OrderStatusRepository) Update(ctx context.Context, id uuid.UUID, patch model.OrderStatusUpdate) (*model.OrderStatusDefinition, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	code := existing.Code
	if patch.Code != nil && strings.TrimSpace(*patch.Code) != "" {
		code = strings.TrimSpace(*patch.Code)
	}
	customer := existing.CustomerLabelRu
	if patch.CustomerLabelRu != nil {
		customer = strings.TrimSpace(*patch.CustomerLabelRu)
		if customer == "" {
			return nil, apperr.BadRequest("customer_label_ru cannot be empty")
		}
	}
	admin := existing.AdminLabelRu
	if patch.AdminLabelRu != nil {
		admin = patch.AdminLabelRu
	}
	desc := existing.Description
	if patch.Description != nil {
		desc = patch.Description
	}
	sortOrder := existing.SortOrder
	if patch.SortOrder != nil {
		sortOrder = *patch.SortOrder
	}
	isActive := existing.IsActive
	if patch.IsActive != nil {
		isActive = *patch.IsActive
	}
	isTerminal := existing.IsTerminal
	if patch.IsTerminal != nil {
		isTerminal = *patch.IsTerminal
	}

	query := `
		UPDATE order_status_definitions SET
			code = $2,
			customer_label_ru = $3,
			admin_label_ru = $4,
			description = $5,
			sort_order = $6,
			is_active = $7,
			is_terminal = $8
		WHERE order_status_id = $1
		RETURNING order_status_id, code, customer_label_ru, admin_label_ru, description,
		          sort_order, is_active, is_terminal, created_at, updated_at
	`
	var d model.OrderStatusDefinition
	err = r.db.Pool.QueryRow(ctx, query, id, code, customer, admin, desc, sortOrder, isActive, isTerminal).Scan(
		&d.OrderStatusID, &d.Code, &d.CustomerLabelRu, &d.AdminLabelRu, &d.Description,
		&d.SortOrder, &d.IsActive, &d.IsTerminal, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperr.Conflict("Order status code already exists")
		}
		return nil, apperr.Internal(err)
	}
	return &d, nil
}

func (r *OrderStatusRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.db.Pool.Exec(ctx, `DELETE FROM order_status_definitions WHERE order_status_id = $1`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return apperr.Conflict("Cannot delete status: orders still use this code")
		}
		return apperr.Internal(err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("%w", apperr.ErrNotFound)
	}
	return nil
}
