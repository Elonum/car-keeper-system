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
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderRepository struct {
	db *database.DB
}

func NewOrderRepository(db *database.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, userID uuid.UUID, create model.OrderCreate, finalPrice float64) (*model.Order, error) {
	var order model.Order
	query := `
		INSERT INTO orders (user_id, configuration_id, status, final_price)
		VALUES ($1, $2, 'pending', $3)
		RETURNING order_id, user_id, configuration_id, manager_id, status, final_price, created_at, updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query, userID, create.ConfigurationID, finalPrice).Scan(
		&order.OrderID, &order.UserID, &order.ConfigurationID, &order.ManagerID,
		&order.Status, &order.FinalPrice, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	return &order, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, orderID uuid.UUID) (*model.OrderWithDetails, error) {
	var order model.OrderWithDetails
	query := `
		SELECT 
			o.order_id, o.user_id, o.configuration_id, o.manager_id, o.status,
			COALESCE(osd.customer_label_ru, o.status) AS status_label,
			o.final_price,
			o.created_at, o.updated_at,
			u.first_name || ' ' || u.last_name as manager_name
		FROM orders o
		LEFT JOIN order_status_definitions osd ON o.status = osd.code
		LEFT JOIN users u ON o.manager_id = u.user_id
		WHERE o.order_id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, orderID).Scan(
		&order.OrderID, &order.UserID, &order.ConfigurationID, &order.ManagerID,
		&order.Status, &order.StatusLabel, &order.FinalPrice, &order.CreatedAt, &order.UpdatedAt,
		&order.ManagerName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w", apperr.ErrNotFound)
		}
		return nil, apperr.Internal(err)
	}

	// Get configuration details
	configRepo := NewConfigurationRepository(r.db)
	config, err := configRepo.GetByID(ctx, order.ConfigurationID)
	if err != nil {
		return nil, err
	}
	order.Configuration = *config

	return &order, nil
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.OrderWithDetails, error) {
	query := `
		SELECT 
			o.order_id, o.user_id, o.configuration_id, o.manager_id, o.status,
			COALESCE(osd.customer_label_ru, o.status) AS status_label,
			o.final_price,
			o.created_at, o.updated_at,
			u.first_name || ' ' || u.last_name as manager_name
		FROM orders o
		LEFT JOIN order_status_definitions osd ON o.status = osd.code
		LEFT JOIN users u ON o.manager_id = u.user_id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	defer rows.Close()

	var orders []model.OrderWithDetails
	for rows.Next() {
		var order model.OrderWithDetails
		if err := rows.Scan(
			&order.OrderID, &order.UserID, &order.ConfigurationID, &order.ManagerID,
			&order.Status, &order.StatusLabel, &order.FinalPrice, &order.CreatedAt, &order.UpdatedAt,
			&order.ManagerName,
		); err != nil {
			return nil, apperr.Internal(err)
		}

		// Get configuration details
		configRepo := NewConfigurationRepository(r.db)
		config, err := configRepo.GetByID(ctx, order.ConfigurationID)
		if err == nil {
			order.Configuration = *config
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// ListAllWithDetails returns all orders with configuration details (staff / admin).
func (r *OrderRepository) ListAllWithDetails(ctx context.Context) ([]model.OrderWithDetails, error) {
	query := `
		SELECT 
			o.order_id, o.user_id, o.configuration_id, o.manager_id, o.status,
			COALESCE(osd.customer_label_ru, o.status) AS status_label,
			o.final_price,
			o.created_at, o.updated_at,
			u.first_name || ' ' || u.last_name as manager_name,
			cust.email::text,
			TRIM(BOTH FROM cust.first_name || ' ' || cust.last_name)::text
		FROM orders o
		LEFT JOIN order_status_definitions osd ON o.status = osd.code
		LEFT JOIN users u ON o.manager_id = u.user_id
		JOIN users cust ON o.user_id = cust.user_id
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	defer rows.Close()

	var orders []model.OrderWithDetails
	for rows.Next() {
		var order model.OrderWithDetails
		if err := rows.Scan(
			&order.OrderID, &order.UserID, &order.ConfigurationID, &order.ManagerID,
			&order.Status, &order.StatusLabel, &order.FinalPrice, &order.CreatedAt, &order.UpdatedAt,
			&order.ManagerName,
			&order.CustomerEmail, &order.CustomerName,
		); err != nil {
			return nil, apperr.Internal(err)
		}

		configRepo := NewConfigurationRepository(r.db)
		config, err := configRepo.GetByID(ctx, order.ConfigurationID)
		if err == nil {
			order.Configuration = *config
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	query := `UPDATE orders SET status = $1 WHERE order_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, orderID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return apperr.BadRequest("Invalid order status")
		}
		return apperr.Internal(err)
	}
	return nil
}

