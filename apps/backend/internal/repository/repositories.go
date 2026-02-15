package repository

import "github.com/pulkitbhatt/ikiru/internal/server"

type Repositories struct {
	UserRepo    *UserRepo
	MonitorRepo *MonitorRepo
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		UserRepo:    NewUserRepo(s.Db.Pool),
		MonitorRepo: NewMonitorRepo(s.Db.Pool),
	}
}
