package repository

import "github.com/pulkitbhatt/ikiru/internal/server"

type Repositories struct {
	UserRepo         *UserRepo
	MonitorRepo      *MonitorRepo
	MonitorCheckRepo *MonitorCheckRepo
	IncidentRepo     *IncidentRepo
	OutboxRepo       *OutboxRepo
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		UserRepo:         NewUserRepo(s.Db.Pool),
		MonitorRepo:      NewMonitorRepo(s.Db.Pool),
		MonitorCheckRepo: NewMonitorCheckRepo(s.Db.Pool),
		IncidentRepo:     NewIncidentRepo(s.Db.Pool, NewOutboxRepo(s.Db.Pool)),
	}
}
