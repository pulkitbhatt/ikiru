package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulkitbhatt/ikiru/internal/model"
)

type MonitorRepo struct {
	db *pgxpool.Pool
}

func NewMonitorRepo(db *pgxpool.Pool) *MonitorRepo {
	return &MonitorRepo{
		db: db,
	}
}

func (mr *MonitorRepo) CreateMonitor(
	ctx context.Context,
	m *model.Monitor,
) error {

	query := `
		INSERT INTO monitors (
			id,
			owner_user_id,
			name,
			description,
			type,
			url,
			interval_seconds,
			timeout_ms,
			status,
			next_check_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`

	_, err := mr.db.Exec(ctx, query,
		m.ID,
		m.OwnerUserID,
		m.Name,
		m.Description,
		m.Type,
		m.URL,
		m.IntervalSeconds,
		m.TimeoutMs,
		m.Status,
		m.NextCheckAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create monitor: %w", err)
	}

	return nil
}
