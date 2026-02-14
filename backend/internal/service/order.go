package service

import (
	"context"
	"fmt"

	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/google/uuid"
)

type OrderService struct {
	repo *repository.Repository
}

func NewOrderService(repos *repository.Repository) *OrderService {
	return &OrderService{repo: repos}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uuid.UUID, create model.OrderCreate) (*model.OrderWithDetails, error) {
	// Get configuration to get final price
	config, err := s.repo.Configuration.GetByID(ctx, create.ConfigurationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	// Verify configuration belongs to user
	if config.UserID != userID {
		return nil, fmt.Errorf("configuration does not belong to user")
	}

	// Verify configuration is in valid status
	if config.Status != "confirmed" && config.Status != "draft" {
		return nil, fmt.Errorf("configuration must be in 'draft' or 'confirmed' status")
	}

	// Create order
	order, err := s.repo.Order.Create(ctx, userID, create, config.TotalPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Update configuration status to 'ordered'
	if err := s.repo.Configuration.UpdateStatus(ctx, create.ConfigurationID, "ordered"); err != nil {
		return nil, fmt.Errorf("failed to update configuration status: %w", err)
	}

	// Get full order details
	orderWithDetails, err := s.repo.Order.GetByID(ctx, order.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order details: %w", err)
	}

	return orderWithDetails, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (*model.OrderWithDetails, error) {
	return s.repo.Order.GetByID(ctx, orderID)
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]model.OrderWithDetails, error) {
	return s.repo.Order.GetByUserID(ctx, userID)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	validStatuses := map[string]bool{
		"pending": true, "approved": true, "paid": true, "completed": true, "cancelled": true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	return s.repo.Order.UpdateStatus(ctx, orderID, status)
}

