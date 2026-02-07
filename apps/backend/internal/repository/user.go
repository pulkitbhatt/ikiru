package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pulkitbhatt/ikiru/internal/server"
	"github.com/pulkitbhatt/ikiru/internal/util"
)

type UserRepo struct {
	server *server.Server
}

func NewUserRepo(s *server.Server) *UserRepo {
	return &UserRepo{
		server: s,
	}
}

func (ur *UserRepo) EnsureUser(ctx context.Context, idpUserId string, email string) (uuid.UUID, error) {
	var userId uuid.UUID

	err := ur.server.Db.Pool.QueryRow(ctx, `
		INSERT INTO users (id, idp_user_id, email)
		VALUES ($1, $2, $3)
		ON CONFLICT (idp_user_id)
		DO UPDATE SET email = COALESCE(users.email, EXCLUDED.email)
		RETURNING id
	`,
		util.GenerateUUID(),
		idpUserId,
		email).Scan(&userId)

	return userId, err
}
