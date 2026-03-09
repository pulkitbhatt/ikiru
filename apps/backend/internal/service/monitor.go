package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pulkitbhatt/ikiru/internal/model"
	"github.com/pulkitbhatt/ikiru/internal/server"
	"github.com/rs/zerolog"
)

type MonitorRepository interface {
	CreateMonitor(ctx context.Context, m *model.Monitor) error
	GetMonitors(ctx context.Context, userID string, limit, offset int) ([]model.Monitor, error)
	GetMonitorById(ctx context.Context, id string, ownerUserID string) (model.Monitor, error)
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

func (s *MonitorService) GetMonitors(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Monitor, error) {
	logger := zerolog.Ctx(ctx)

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	monitors, err := s.monitorRepo.GetMonitors(ctx, userID.String(), limit, offset)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID.String()).Msg("failed to get monitors")
		return nil, fmt.Errorf("failed to get monitors: %w", err)
	}

	return monitors, nil
}

func (s *MonitorService) GetMonitorById(ctx context.Context, userID, id uuid.UUID) (model.Monitor, error) {
	logger := zerolog.Ctx(ctx)
	monitor, err := s.monitorRepo.GetMonitorById(ctx, id.String(), userID.String())
	if err != nil {
		logger.Error().
			Err(err).
			Str("id", id.String()).
			Str("user_id", userID.String()).
			Msg("failed to get monitor with given id")
		return model.Monitor{}, fmt.Errorf("failed to get monitor: %w", err)
	}
	return monitor, nil
}