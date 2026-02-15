package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pulkitbhatt/ikiru/internal/handler/dto"
	"github.com/pulkitbhatt/ikiru/internal/model"
	"github.com/pulkitbhatt/ikiru/internal/server"
	"github.com/pulkitbhatt/ikiru/internal/service"
	"github.com/pulkitbhatt/ikiru/internal/validation"
)

type MonitorHandler struct {
	Handler
	monitorService *service.MonitorService
}

func NewMonitorHandler(server *server.Server, monitorService *service.MonitorService) *MonitorHandler {
	return &MonitorHandler{
		Handler:        *NewHandler(server),
		monitorService: monitorService,
	}
}

func (h *MonitorHandler) CreateMonitor(c echo.Context) error {
	ctx := c.Request().Context()
	log := LoggerFromContext(c)

	userID := GetUserID(c)

	var req dto.CreateMonitorRequest
	if err := c.Bind(&req); err != nil {
		log.Warn().Err(err).Msg("failed to bind create monitor request")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	monitor := model.NewMonitor(
		userID,
		req.Name,
		req.URL,
		req.IntervalSeconds,
		req.TimeoutMs,
		req.Description,
	)

	if err := h.monitorService.CreateMonitor(ctx, monitor); err != nil {
		var verr validation.ValidationErrors
		if errors.As(err, &verr) {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"success": false,
				"errors":  verr,
			})
		}

		log.Error().Err(err).Msg("failed to create monitor")
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create monitor")
	}

	log.Info().
		Str("monitor_id", monitor.ID.String()).
		Msg("monitor created")

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"data": map[string]any{
			"id": monitor.ID,
		},
	})
}
