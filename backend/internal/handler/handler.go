package handler

import (
	"github.com/carkeeper/backend/config"
	"github.com/carkeeper/backend/internal/service"
)

type Handler struct {
	services *service.Service
	cfg      *config.Config
}

func New(services *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
	}
}

func (h *Handler) Services() *service.Service {
	return h.services
}


