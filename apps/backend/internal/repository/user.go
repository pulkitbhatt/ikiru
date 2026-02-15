package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulkitbhatt/ikiru/internal/util"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) EnsureUser(ctx context.Context, idpUserId string, email string) (uuid.UUID, error) {
	var userId uuid.UUID

	query := `
		INSERT INTO users (id, idp_user_id, email)
		VALUES ($1, $2, $3)
		ON CONFLICT (idp_user_id)
		DO UPDATE SET email = COALESCE(users.email, EXCLUDED.email)
		RETURNING id
	`

	err := ur.db.QueryRow(ctx, query,
		util.GenerateUUIDStr(),
		idpUserId,
		email,
	).Scan(&userId)

	return userId, err
}
