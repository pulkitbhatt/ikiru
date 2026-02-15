package service

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	"github.com/pulkitbhatt/ikiru/internal/server"
)

type UserRepository interface {
	EnsureUser(ctx context.Context, idpUserID string, email string) (uuid.UUID, error)
}

type AuthService struct {
	server   *server.Server
	userRepo UserRepository
}

func NewAuthService(s *server.Server, userRepo UserRepository) *AuthService {
	clerk.SetKey(s.Config.Auth.SecretKey)
	return &AuthService{
		server:   s,
		userRepo: userRepo,
	}
}

func (u *AuthService) EnsureUser(ctx context.Context, idpUserID string, email string) (uuid.UUID, error) {
	return u.userRepo.EnsureUser(ctx, idpUserID, email)
}
