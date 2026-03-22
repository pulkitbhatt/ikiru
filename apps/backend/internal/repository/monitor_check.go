package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulkitbhatt/ikiru/internal/model"
)

type MonitorCheckRepo struct {
	db *pgxpool.Pool
}

func NewMonitorCheckRepo(db *pgxpool.Pool) *MonitorCheckRepo {
	return &MonitorCheckRepo{
		db: db,
	}
}

func (r *MonitorCheckRepo) InsertCheckResult(
	ctx context.Context,
	rec model.MonitorCheckResult,
) error {

	query := `
        INSERT INTO monitor_check_results (
            id, monitor_id, region, scheduled_at,
            started_at, finished_at, status,
            http_status, latency_ms, error
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
        ON CONFLICT (monitor_id, region, scheduled_at)
        DO NOTHING;
    `

	_, err := r.db.Exec(ctx, query,
		rec.ID,
		rec.MonitorID,
		rec.Region,
		rec.ScheduledAt,
		rec.StartedAt,
		rec.FinishedAt,
		rec.Status,
		rec.HTTPStatus,
		rec.LatencyMs,
		rec.Error,
	)
	return err
}

func (r *MonitorCheckRepo) GetLastNResults(
	ctx context.Context,
	monitorID uuid.UUID,
	region string,
	n int,
) ([]string, error) {

	query := `
		SELECT status
		FROM monitor_check_results
		WHERE monitor_id = $1 AND region = $2
		ORDER BY scheduled_at DESC
		LIMIT $3
	`
	rows, err := r.db.Query(ctx, query, monitorID, region, n)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var s string
		rows.Scan(&s)
		results = append(results, s)
	}

	return results, nil
}
