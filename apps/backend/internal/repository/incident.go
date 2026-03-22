package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulkitbhatt/ikiru/internal/model"
	"github.com/pulkitbhatt/ikiru/internal/util"
)

type IncidentRepo struct {
	db         *pgxpool.Pool
	outboxRepo *OutboxRepo
}

func NewIncidentRepo(db *pgxpool.Pool, outbox *OutboxRepo) *IncidentRepo {
	return &IncidentRepo{
		db:         db,
		outboxRepo: outbox,
	}
}

func (r *IncidentRepo) GetOpenIncident(
	ctx context.Context,
	monitorID uuid.UUID,
	region string,
) (*model.Incident, error) {

	query := `
		SELECT id, started_at
		FROM incidents
		WHERE monitor_id = $1 AND region = $2 AND status = 'open'
		LIMIT 1
	`
	row := r.db.QueryRow(ctx, query, monitorID, region)

	var i model.Incident
	err := row.Scan(&i.ID, &i.StartedAt)
	if err != nil {
		return nil, nil
	}

	return &i, nil
}

func (r *IncidentRepo) TryCreateIncidentWithOutbox(
	ctx context.Context,
	monitorID uuid.UUID,
	region string,
	failureCount int,
) (bool, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO incidents (id, monitor_id, region, status, started_at, failure_count)
		VALUES ($1, $2, $3, 'open', now(), $4)
		ON CONFLICT DO NOTHING
	`

	incidentID := util.GenerateUUID()
	res, err := tx.Exec(ctx, query,
		incidentID,
		monitorID,
		region,
		failureCount,
	)
	if err != nil {
		return false, err
	}

	if res.RowsAffected() == 0 {
		return false, nil
	}

	event := map[string]any{
		"incident_id": incidentID.String(),
		"monitor_id":  monitorID.String(),
		"region":      region,
		"type":        "created",
		"timestamp":   time.Now().UTC(),
	}

	payload, _ := json.Marshal(event)

	if err := r.outboxRepo.InsertOutboxEventTx(
		ctx,
		tx,
		"incident.created",
		payload,
	); err != nil {
		return false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *IncidentRepo) ResolveIncidentWithOutbox(
	ctx context.Context,
	monitorID uuid.UUID,
	region string,
) (bool, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	var incidentID uuid.UUID

	query := `
		UPDATE incidents
		SET status = 'resolved',
			resolved_at = now(),
			updated_at = now()
		WHERE monitor_id = $1
		AND region = $2
		AND status = 'open'
		RETURNING id
	`
	err = tx.QueryRow(ctx, query,
		monitorID,
		region,
	).Scan(&incidentID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	event := map[string]any{
		"incident_id": incidentID.String(),
		"monitor_id":  monitorID.String(),
		"region":      region,
		"type":        "resolved",
		"timestamp":   time.Now().UTC(),
	}

	payload, _ := json.Marshal(event)

	if err := r.outboxRepo.InsertOutboxEventTx(
		ctx,
		tx,
		"incident.resolved",
		payload,
	); err != nil {
		return false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *IncidentRepo) IncrementFailureCount(
	ctx context.Context,
	monitorID uuid.UUID,
	region string,
) error {
	query := `
		UPDATE incidents
		SET failure_count = failure_count + 1,
			updated_at = now()
		WHERE monitor_id = $1 AND region = $2 AND status = 'open'
	`
	_, err := r.db.Exec(ctx, query, monitorID, region)
	return err
}
