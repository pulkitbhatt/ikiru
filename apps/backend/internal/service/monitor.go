package service

import (
	"context"
	"fmt"

	"github.com/pulkitbhatt/ikiru/internal/model"
	"github.com/pulkitbhatt/ikiru/internal/server"
	"github.com/rs/zerolog"
)

type MonitorRepository interface {
	CreateMonitor(ctx context.Context, m *model.Monitor) error
}

type MonitorService struct {
	server      *server.Server
	monitorRepo MonitorRepository
}

func NewMonitorService(server *server.Server, monitorRepo MonitorRepository) *MonitorService {
	return &MonitorService{
		server:      server,
		monitorRepo: monitorRepo,
	}
}

func (s *MonitorService) CreateMonitor(ctx context.Context, m *model.Monitor) error {
	logger := zerolog.Ctx(ctx)
	if err := m.Validate(); err != nil {
		logger.Debug().Err(err).Msg("validation failed for creating monitor")
		return err
	}

	if err := s.monitorRepo.CreateMonitor(ctx, m); err != nil {
		logger.Error().Err(err).Str("monitor_id", m.ID.String()).Msg("failed to create monitor")
		return fmt.Errorf("failed to create monitor: %w", err)
	}

	return nil
}
