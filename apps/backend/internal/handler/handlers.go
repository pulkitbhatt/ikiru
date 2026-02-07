package handler

import (
	"github.com/pulkitbhatt/ikiru/internal/server"
	"github.com/pulkitbhatt/ikiru/internal/service"
)

type Handlers struct {
	Health   *HealthHandler
	Monitor  *MonitorHandler
	Services *service.Services
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(s),
		Monitor:  NewMonitorHandler(s),
		Services: services,
	}
}
