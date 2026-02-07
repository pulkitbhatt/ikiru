package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pulkitbhatt/ikiru/internal/server"
	"github.com/rs/zerolog"
)

const (
	UserIDKey = "user_id"
	LoggerKey = "logger"
)

type Handler struct {
	server *server.Server
}

func NewHandler(s *server.Server) *Handler {
	return &Handler{
		server: s,
	}
}

func GetUserID(c echo.Context) uuid.UUID {
	if id, ok := c.Get(UserIDKey).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

func LoggerFromContext(c echo.Context) zerolog.Logger {
	if l, ok := c.Get(LoggerKey).(*zerolog.Logger); ok {
		return *l
	}
	return zerolog.Nop()
}
