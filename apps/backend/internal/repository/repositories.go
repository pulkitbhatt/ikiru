package repository

import "github.com/pulkitbhatt/ikiru/internal/server"

type Repositories struct {
	UserRepo *UserRepo
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		UserRepo: NewUserRepo(s),
	}
}
