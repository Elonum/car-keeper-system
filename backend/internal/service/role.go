package service

import (
	"context"

	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
)

type RoleService struct {
	repo *repository.Repository
}

func NewRoleService(repos *repository.Repository) *RoleService {
	return &RoleService{repo: repos}
}

func (s *RoleService) ListDefinitions(ctx context.Context) ([]model.RoleDefinition, error) {
	return s.repo.Role.ListAll(ctx)
}
