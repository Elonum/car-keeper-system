package service

import (
	"context"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/validate"
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
	code, msg := validate.OrderStatusCode(in.Code)
	if msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	customer, msg := validate.OrderStatusCustomerLabel(in.CustomerLabelRu)
	if msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	admin, msg := validate.OrderStatusAdminLabel(in.AdminLabelRu)
	if msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	desc, msg := validate.OrderStatusDescription(in.Description)
	if msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	sortOrder := 0
	if in.SortOrder != nil {
		sortOrder = *in.SortOrder
	}
	if msg := validate.OrderStatusSortOrder(sortOrder); msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	in.Code = code
	in.CustomerLabelRu = customer
	in.AdminLabelRu = admin
	in.Description = desc
	return s.repo.OrderStatus.Create(ctx, in)
}

func (s *OrderStatusService) Update(ctx context.Context, id uuid.UUID, patch model.OrderStatusUpdate) (*model.OrderStatusDefinition, error) {
	if patch.Code != nil {
		code, msg := validate.OrderStatusCode(*patch.Code)
		if msg != "" {
			return nil, apperr.BadRequest(msg)
		}
		patch.Code = &code
	}
	if patch.CustomerLabelRu != nil {
		customer, msg := validate.OrderStatusCustomerLabel(*patch.CustomerLabelRu)
		if msg != "" {
			return nil, apperr.BadRequest(msg)
		}
		patch.CustomerLabelRu = &customer
	}
	if patch.AdminLabelRu != nil {
		admin, msg := validate.OrderStatusAdminLabel(patch.AdminLabelRu)
		if msg != "" {
			return nil, apperr.BadRequest(msg)
		}
		patch.AdminLabelRu = admin
	}
	if patch.Description != nil {
		desc, msg := validate.OrderStatusDescription(patch.Description)
		if msg != "" {
			return nil, apperr.BadRequest(msg)
		}
		patch.Description = desc
	}
	if patch.SortOrder != nil {
		if msg := validate.OrderStatusSortOrder(*patch.SortOrder); msg != "" {
			return nil, apperr.BadRequest(msg)
		}
	}
	return s.repo.OrderStatus.Update(ctx, id, patch)
}

func (s *OrderStatusService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.OrderStatus.Delete(ctx, id)
}
