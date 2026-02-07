package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pulkitbhatt/ikiru/internal/server"
)

type HealthHandler struct {
	Handler
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{
		Handler: *NewHandler(s),
	}
}

func (h *HealthHandler) CheckHealth(c echo.Context) error {
	c.JSON(http.StatusOK, map[string]string{
		"message": "OK",
	})
	return nil
}
