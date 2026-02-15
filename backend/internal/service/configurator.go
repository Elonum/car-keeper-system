package service

import (
	"context"
	"fmt"
	"log"

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
	// Get trim to calculate base price
	trim, err := s.repo.Trim.GetByID(ctx, create.TrimID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trim: %w", err)
	}

	// Get color to get price delta
	color, err := s.repo.Color.GetByID(ctx, create.ColorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get color: %w", err)
	}

	// Calculate total price
	totalPrice := trim.BasePrice + color.PriceDelta

	// Get selected options and add their prices
	if len(create.OptionIDs) > 0 {
		options, err := s.repo.Option.GetByIDs(ctx, create.OptionIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get options: %w", err)
		}
		for _, opt := range options {
			totalPrice += opt.Price
		}
	}

	// Create configuration
	config, err := s.repo.Configuration.Create(ctx, userID, create, totalPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to create configuration: %w", err)
	}

	// Get full details
	configWithDetails, err := s.repo.Configuration.GetByID(ctx, config.ConfigurationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration details: %w", err)
	}

	return configWithDetails, nil
}

func (s *ConfiguratorService) GetConfiguration(ctx context.Context, configID uuid.UUID) (*model.ConfigurationWithDetails, error) {
	return s.repo.Configuration.GetByID(ctx, configID)
}

func (s *ConfiguratorService) GetUserConfigurations(ctx context.Context, userID uuid.UUID) ([]model.ConfigurationWithDetails, error) {
	return s.repo.Configuration.GetByUserID(ctx, userID)
}

func (s *ConfiguratorService) UpdateConfigurationStatus(ctx context.Context, configID uuid.UUID, status string) error {
	log.Printf("[UpdateConfigurationStatus] Called with configID: %s, status: '%s'", configID.String(), status)
	
	if status == "" {
		log.Printf("[UpdateConfigurationStatus] Error: status is empty")
		return fmt.Errorf("status cannot be empty")
	}
	
	validStatuses := map[string]bool{
		"draft": true, "confirmed": true, "ordered": true, "cancelled": true, "purchased": true,
	}
	if !validStatuses[status] {
		log.Printf("[UpdateConfigurationStatus] Error: invalid status '%s'", status)
		return fmt.Errorf("invalid status: %s", status)
	}

	log.Printf("[UpdateConfigurationStatus] Updating status to '%s'", status)
	return s.repo.Configuration.UpdateStatus(ctx, configID, status)
}

func (s *ConfiguratorService) UpdateConfiguration(ctx context.Context, configID uuid.UUID, userID uuid.UUID, update model.ConfigurationCreate) (*model.ConfigurationWithDetails, error) {
	// Verify configuration belongs to user
	existingConfig, err := s.repo.Configuration.GetByID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("configuration not found: %w", err)
	}
	if existingConfig.UserID != userID {
		return nil, fmt.Errorf("unauthorized: configuration does not belong to user")
	}

	// Only allow updating draft configurations
	if existingConfig.Status != "draft" {
		return nil, fmt.Errorf("can only update draft configurations")
	}

	// Get trim to calculate base price
	trim, err := s.repo.Trim.GetByID(ctx, update.TrimID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trim: %w", err)
	}

	// Get color to get price delta
	color, err := s.repo.Color.GetByID(ctx, update.ColorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get color: %w", err)
	}

	// Calculate total price
	totalPrice := trim.BasePrice + color.PriceDelta

	// Get selected options and add their prices
	if len(update.OptionIDs) > 0 {
		options, err := s.repo.Option.GetByIDs(ctx, update.OptionIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get options: %w", err)
		}
		for _, opt := range options {
			totalPrice += opt.Price
		}
	}

	// Update configuration
	updatedConfig, err := s.repo.Configuration.Update(ctx, configID, update, totalPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to update configuration: %w", err)
	}

	// Get full details
	configWithDetails, err := s.repo.Configuration.GetByID(ctx, updatedConfig.ConfigurationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration details: %w", err)
	}

	return configWithDetails, nil
}

func (s *ConfiguratorService) DeleteConfiguration(ctx context.Context, configID uuid.UUID) error {
	// Get configuration to check status
	config, err := s.repo.Configuration.GetByID(ctx, configID)
	if err != nil {
		return fmt.Errorf("configuration not found: %w", err)
	}

	// Only allow deleting draft configurations
	if config.Status != "draft" {
		return fmt.Errorf("can only delete draft configurations, current status: %s", config.Status)
	}

	return s.repo.Configuration.Delete(ctx, configID)
}

