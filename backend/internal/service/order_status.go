package service

import (
	"context"

	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/google/uuid"
)

type OrderStatusService struct {
	repo *repository.Repository
}

func NewOrderStatusService(repos *repository.Repository) *OrderStatusService {
	return &OrderStatusService{repo: repos}
}

func (s *OrderStatusService) ListPublic(ctx context.Context) ([]model.OrderStatusPublic, error) {
	list, err := s.repo.OrderStatus.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]model.OrderStatusPublic, 0, len(list))
	for _, d := range list {
		out = append(out, model.OrderStatusPublic{
			Code:            d.Code,
			CustomerLabelRu: d.CustomerLabelRu,
			IsTerminal:      d.IsTerminal,
			SortOrder:       d.SortOrder,
		})
	}
	return out, nil
}

func (s *OrderStatusService) ListAll(ctx context.Context) ([]model.OrderStatusDefinition, error) {
	return s.repo.OrderStatus.ListAll(ctx)
}

func (s *OrderStatusService) Create(ctx context.Context, in model.OrderStatusCreate) (*model.OrderStatusDefinition, error) {
	return s.repo.OrderStatus.Create(ctx, in)
}

func (s *OrderStatusService) Update(ctx context.Context, id uuid.UUID, patch model.OrderStatusUpdate) (*model.OrderStatusDefinition, error) {
	return s.repo.OrderStatus.Update(ctx, id, patch)
}

func (s *OrderStatusService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.OrderStatus.Delete(ctx, id)
}
