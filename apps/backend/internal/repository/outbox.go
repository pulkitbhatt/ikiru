package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulkitbhatt/ikiru/internal/model"
	"github.com/pulkitbhatt/ikiru/internal/util"
)

type OutboxRepo struct {
	db *pgxpool.Pool
}

func NewOutboxRepo(db *pgxpool.Pool) *OutboxRepo {
	return &OutboxRepo{
		db: db,
	}
}

func (r *OutboxRepo) InsertOutboxEventTx(
	ctx context.Context,
	tx pgx.Tx,
	eventType string,
	payload []byte,
) error {

	query := `
		INSERT INTO outbox_events (id, event_type, payload)
		VALUES ($1, $2, $3)
	`
	_, err := tx.Exec(ctx, query,
		util.GenerateUUID(),
		eventType,
		payload,
	)

	return err
}

func (r *OutboxRepo) FetchUnprocessed(
	ctx context.Context,
	limit int,
) ([]model.OutboxEvent, error) {

	query := `
		UPDATE outbox_events
		SET processing_at = now()
		WHERE id IN (
			SELECT id
			FROM outbox_events
			WHERE processed_at IS NULL
			AND (
				processing_at IS NULL
				OR processing_at < now() - interval '1 minute'
			)
			ORDER BY created_at
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, event_type, payload;
	`
	rows, err := r.db.Query(ctx, query, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.OutboxEvent

	for rows.Next() {
		var e model.OutboxEvent
		if err := rows.Scan(&e.ID, &e.Type, &e.Payload); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (r *OutboxRepo) MarkProcessed(
	ctx context.Context,
	id uuid.UUID,
) error {

	query := `
		UPDATE outbox_events
		SET processed_at = now()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)

	return err
}
