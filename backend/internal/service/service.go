package service

import (
	"github.com/carkeeper/backend/config"
	"github.com/carkeeper/backend/internal/repository"
)

type Service struct {
	Auth        *AuthService
	Catalog     *CatalogService
	Configurator *ConfiguratorService
	Order       *OrderService
	Service     *ServiceService
	News        *NewsService
	Profile     *ProfileService
}

func New(repos *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		Auth:        NewAuthService(repos, cfg),
		Catalog:     NewCatalogService(repos),
		Configurator: NewConfiguratorService(repos),
		Order:       NewOrderService(repos),
		Service:     NewServiceService(repos),
		News:        NewNewsService(repos),
		Profile:     NewProfileService(repos),
	}
}

