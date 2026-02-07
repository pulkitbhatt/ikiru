package service

import (
	"github.com/pulkitbhatt/ikiru/internal/repository"
	"github.com/pulkitbhatt/ikiru/internal/server"
)

type UserService struct {
	server   *server.Server
	userRepo *repository.UserRepo
}
