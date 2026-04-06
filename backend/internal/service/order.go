package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/authz"
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
	config, err := s.repo.Configuration.GetByID(ctx, create.ConfigurationID)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, apperr.NotFoundErr("Configuration not found")
		}
		return nil, err
	}

	if config.UserID != userID {
		return nil, apperr.Forbidden("Configuration does not belong to your account")
	}

	if config.Status != "confirmed" && config.Status != "draft" {
		return nil, apperr.BadRequest("Configuration cannot be ordered in its current status")
	}

	order, err := s.repo.Order.Create(ctx, userID, create, config.TotalPrice)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Configuration.UpdateStatus(ctx, create.ConfigurationID, "ordered"); err != nil {
		return nil, err
	}

	orderWithDetails, err := s.repo.Order.GetByID(ctx, order.OrderID)
	if err != nil {
		return nil, err
	}

	return orderWithDetails, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID, requester uuid.UUID, role string) (*model.OrderWithDetails, error) {
	order, err := s.repo.Order.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if !authz.CanViewOrder(order.UserID, requester, role) {
		return nil, fmt.Errorf("%w", apperr.ErrNotFound)
	}
	return order, nil
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]model.OrderWithDetails, error) {
	return s.repo.Order.GetByUserID(ctx, userID)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string, requester uuid.UUID, role string) error {
	status = strings.TrimSpace(status)
	if status == "" {
		return apperr.BadRequest("status is required")
	}

	target, err := s.repo.OrderStatus.GetByCode(ctx, status)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return apperr.BadRequest("Unknown order status")
		}
		return err
	}

	if !authz.HasPermission(role, authz.PermOrdersManageStatus) && !target.IsActive {
		return apperr.BadRequest("This order status is not available")
	}

	order, err := s.repo.Order.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	if !authz.CanUpdateOrderStatus(order.UserID, requester, role, order.Status, status) {
		return fmt.Errorf("%w", apperr.ErrForbidden)
	}

	return s.repo.Order.UpdateStatus(ctx, orderID, status)
}
