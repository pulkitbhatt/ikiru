package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulkitbhatt/ikiru/internal/handler/dto"
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

func (r *MonitorRepo) ClaimDueMonitors(
	ctx context.Context,
	limit int,
) ([]dto.DueMonitor, error) {

	const query = `
UPDATE monitors
SET next_check_at = GREATEST(
    next_check_at + (interval_seconds * INTERVAL '1 second'),
    now()
)
WHERE id IN (
    SELECT id
    FROM monitors
    WHERE
        next_check_at <= now()
        AND status = 'active'
        AND deleted_at IS NULL
    ORDER BY next_check_at
    FOR UPDATE SKIP LOCKED
    LIMIT $1
)
RETURNING
    id,
    owner_user_id,
    type,
    url,
    timeout_ms
`

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dueMonitors := make([]dto.DueMonitor, 0, limit)

	for rows.Next() {
		var m dto.DueMonitor

		err := rows.Scan(
			&m.ID,
			&m.UserId,
			&m.Type,
			&m.URL,
			&m.TimeoutMs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan the current row into Monitor struct: %w", err)
		}

		dueMonitors = append(dueMonitors, m)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred while claiming monitors: %w", rows.Err())
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error occurred while committing the transaction: %w", err)
	}

	return dueMonitors, nil
}
