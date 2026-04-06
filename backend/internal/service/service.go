package service

import (
	"github.com/carkeeper/backend/config"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/storage"
)

type Service struct {
	Auth          *AuthService
	Catalog       *CatalogService
	Configurator  *ConfiguratorService
	Order         *OrderService
	OrderStatus   *OrderStatusService
	Role          *RoleService
	Service       *ServiceService
	News          *NewsService
	Profile       *ProfileService
	Document      *DocumentService
}

func New(repos *repository.Repository, cfg *config.Config, fileStore storage.FileStorage) *Service {
	return &Service{
		Auth:         NewAuthService(repos, cfg),
		Catalog:      NewCatalogService(repos),
		Configurator: NewConfiguratorService(repos),
		Order:        NewOrderService(repos),
		OrderStatus:  NewOrderStatusService(repos),
		Role:         NewRoleService(repos),
		Service:      NewServiceService(repos),
		News:         NewNewsService(repos),
		Profile:      NewProfileService(repos),
		Document:     NewDocumentService(repos, fileStore, cfg.Storage.MaxUploadBytes),
	}
}

