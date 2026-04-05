package service

import (
	"context"
	"fmt"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/authz"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/google/uuid"
)

type ConfiguratorService struct {
	repo *repository.Repository
}

func NewConfiguratorService(repos *repository.Repository) *ConfiguratorService {
	return &ConfiguratorService{repo: repos}
}

func (s *ConfiguratorService) GetColors(ctx context.Context, isAvailable *bool) ([]model.Color, error) {
	return s.repo.Color.GetAll(ctx, isAvailable)
}

func (s *ConfiguratorService) GetOptions(ctx context.Context, trimID uuid.UUID) ([]model.Option, error) {
	return s.repo.Option.GetByTrimID(ctx, trimID)
}

func (s *ConfiguratorService) CreateConfiguration(ctx context.Context, userID uuid.UUID, create model.ConfigurationCreate) (*model.ConfigurationWithDetails, error) {
	trim, err := s.repo.Trim.GetByID(ctx, create.TrimID)
	if err != nil {
		return nil, apperr.Wrap(err, 400, "Invalid trim or catalog data")
	}

	color, err := s.repo.Color.GetByID(ctx, create.ColorID)
	if err != nil {
		return nil, apperr.Wrap(err, 400, "Invalid color or catalog data")
	}

	totalPrice := trim.BasePrice + color.PriceDelta

	if len(create.OptionIDs) > 0 {
		options, err := s.repo.Option.GetByIDs(ctx, create.OptionIDs)
		if err != nil {
			return nil, apperr.Wrap(err, 400, "Invalid options selection")
		}
		for _, opt := range options {
			totalPrice += opt.Price
		}
	}

	config, err := s.repo.Configuration.Create(ctx, userID, create, totalPrice)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	configWithDetails, err := s.repo.Configuration.GetByID(ctx, config.ConfigurationID)
	if err != nil {
		return nil, err
	}

	return configWithDetails, nil
}

func (s *ConfiguratorService) GetConfiguration(ctx context.Context, configID uuid.UUID, requester uuid.UUID, role string) (*model.ConfigurationWithDetails, error) {
	config, err := s.repo.Configuration.GetByID(ctx, configID)
	if err != nil {
		return nil, err
	}
	if !authz.CanAccessConfiguration(config.UserID, requester, role) {
		return nil, fmt.Errorf("%w", apperr.ErrNotFound)
	}
	return config, nil
}

func (s *ConfiguratorService) GetUserConfigurations(ctx context.Context, userID uuid.UUID) ([]model.ConfigurationWithDetails, error) {
	return s.repo.Configuration.GetByUserID(ctx, userID)
}

func (s *ConfiguratorService) UpdateConfigurationStatus(ctx context.Context, configID uuid.UUID, requester uuid.UUID, role, status string) error {
	if status == "" {
		return apperr.BadRequest("status cannot be empty")
	}

	validStatuses := map[string]bool{
		"draft": true, "confirmed": true, "ordered": true, "cancelled": true, "purchased": true,
	}
	if !validStatuses[status] {
		return apperr.BadRequest("Invalid configuration status")
	}

	config, err := s.repo.Configuration.GetByID(ctx, configID)
	if err != nil {
		return err
	}

	if authz.IsStaff(role) {
		return s.repo.Configuration.UpdateStatus(ctx, configID, status)
	}

	if config.UserID != requester {
		return fmt.Errorf("%w", apperr.ErrForbidden)
	}
	if !authz.CustomerMayChangeConfigurationStatus(config.Status, status) {
		return fmt.Errorf("%w", apperr.ErrForbidden)
	}

	return s.repo.Configuration.UpdateStatus(ctx, configID, status)
}

func (s *ConfiguratorService) UpdateConfiguration(ctx context.Context, configID uuid.UUID, userID uuid.UUID, role string, update model.ConfigurationCreate) (*model.ConfigurationWithDetails, error) {
	existingConfig, err := s.repo.Configuration.GetByID(ctx, configID)
	if err != nil {
		return nil, err
	}
	if !authz.CanAccessConfiguration(existingConfig.UserID, userID, role) {
		return nil, fmt.Errorf("%w", apperr.ErrNotFound)
	}

	if existingConfig.Status != "draft" {
		return nil, apperr.BadRequest("Only draft configurations can be updated")
	}

	trim, err := s.repo.Trim.GetByID(ctx, update.TrimID)
	if err != nil {
		return nil, apperr.Wrap(err, 400, "Invalid trim or catalog data")
	}

	color, err := s.repo.Color.GetByID(ctx, update.ColorID)
	if err != nil {
		return nil, apperr.Wrap(err, 400, "Invalid color or catalog data")
	}

	totalPrice := trim.BasePrice + color.PriceDelta

	if len(update.OptionIDs) > 0 {
		options, err := s.repo.Option.GetByIDs(ctx, update.OptionIDs)
		if err != nil {
			return nil, apperr.Wrap(err, 400, "Invalid options selection")
		}
		for _, opt := range options {
			totalPrice += opt.Price
		}
	}

	updatedConfig, err := s.repo.Configuration.Update(ctx, configID, update, totalPrice)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	configWithDetails, err := s.repo.Configuration.GetByID(ctx, updatedConfig.ConfigurationID)
	if err != nil {
		return nil, err
	}

	return configWithDetails, nil
}

func (s *ConfiguratorService) DeleteConfiguration(ctx context.Context, configID uuid.UUID, requester uuid.UUID, role string) error {
	config, err := s.repo.Configuration.GetByID(ctx, configID)
	if err != nil {
		return err
	}
	if !authz.CanAccessConfiguration(config.UserID, requester, role) {
		return fmt.Errorf("%w", apperr.ErrNotFound)
	}
	if config.Status != "draft" {
		return apperr.BadRequest("Only draft configurations can be deleted")
	}

	return s.repo.Configuration.Delete(ctx, configID)
}
