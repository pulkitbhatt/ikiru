package service

import (
	"context"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	"github.com/pulkitbhatt/ikiru/internal/repository"
	"github.com/pulkitbhatt/ikiru/internal/server"
)

type AuthService struct {
	server   *server.Server
	userRepo *repository.UserRepo
}

func NewAuthService(s *server.Server) *AuthService {
	clerk.SetKey(s.Config.Auth.SecretKey)
	return &AuthService{
		server:   s,
		userRepo: repository.NewUserRepo(s),
	}
}

func (u *AuthService) EnsureUser(ctx context.Context, idpUserID string, email string) (uuid.UUID, error) {
	return u.userRepo.EnsureUser(ctx, idpUserID, email)
}
