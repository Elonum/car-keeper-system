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
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, orderID uuid.UUID) (*model.OrderWithDetails, error) {
	var order model.OrderWithDetails
	query := `
		SELECT 
			o.order_id, o.user_id, o.configuration_id, o.manager_id, o.status, o.final_price,
			o.created_at, o.updated_at,
			u.first_name || ' ' || u.last_name as manager_name
		FROM orders o
		LEFT JOIN users u ON o.manager_id = u.user_id
		WHERE o.order_id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, orderID).Scan(
		&order.OrderID, &order.UserID, &order.ConfigurationID, &order.ManagerID,
		&order.Status, &order.FinalPrice, &order.CreatedAt, &order.UpdatedAt,
		&order.ManagerName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get configuration details
	configRepo := NewConfigurationRepository(r.db)
	config, err := configRepo.GetByID(ctx, order.ConfigurationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}
	order.Configuration = *config

	return &order, nil
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.OrderWithDetails, error) {
	query := `
		SELECT 
			o.order_id, o.user_id, o.configuration_id, o.manager_id, o.status, o.final_price,
			o.created_at, o.updated_at,
			u.first_name || ' ' || u.last_name as manager_name
		FROM orders o
		LEFT JOIN users u ON o.manager_id = u.user_id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []model.OrderWithDetails
	for rows.Next() {
		var order model.OrderWithDetails
		if err := rows.Scan(
			&order.OrderID, &order.UserID, &order.ConfigurationID, &order.ManagerID,
			&order.Status, &order.FinalPrice, &order.CreatedAt, &order.UpdatedAt,
			&order.ManagerName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
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

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	query := `UPDATE orders SET status = $1 WHERE order_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}

