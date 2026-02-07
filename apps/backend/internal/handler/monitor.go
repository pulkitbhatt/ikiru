package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pulkitbhatt/ikiru/internal/server"
)

type MonitorHandler struct {
	Handler
}

func NewMonitorHandler(s *server.Server) *MonitorHandler {
	return &MonitorHandler{
		Handler: *NewHandler(s),
	}
}

func (m *MonitorHandler) CreateMonitor(c echo.Context) error {
	log := LoggerFromContext(c)
	log.Info().Msgf("user id: %v", GetUserID(c))
	log.Info().Msgf("idp id: %v", c.Get("idp_user_id").(string))
	c.JSON(http.StatusCreated, map[string]any{
		"success": true,
	})
	return nil
}
