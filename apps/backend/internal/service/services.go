package service

import (
	"github.com/pulkitbhatt/ikiru/internal/repository"
	"github.com/pulkitbhatt/ikiru/internal/server"
)

type Services struct {
	Auth           AuthService
	MonitorService MonitorService
}

func NewServices(s *server.Server, repos *repository.Repositories) *Services {
	return &Services{
		Auth:           *NewAuthService(s, repos.UserRepo),
		MonitorService: *NewMonitorService(s, repos.MonitorRepo),
	}
}
