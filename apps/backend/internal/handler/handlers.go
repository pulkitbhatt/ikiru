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

func NewHandlers(server *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(server),
		Monitor:  NewMonitorHandler(server, &services.MonitorService),
		Services: services,
	}
}
