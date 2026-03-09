package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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

func (h *MonitorHandler) GetMonitors(c echo.Context) error {
	ctx := c.Request().Context()
	log := LoggerFromContext(c)
	userId := GetUserID(c)

	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	var (
		limit  int
		offset int
		err    error
	)

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Warn().Str("limit", limitStr).Msg("invalid limit query parameter")
			return echo.NewHTTPError(http.StatusBadRequest, "invalid limit query parameter")
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			log.Warn().Str("offset", offsetStr).Msg("invalid offset query parameter")
			return echo.NewHTTPError(http.StatusBadRequest, "invalid offset query parameter")
		}
	}

	monitors, err := h.monitorService.GetMonitors(ctx, userId, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to get monitors")
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get monitors")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{"monitors": monitors},
	})
}

func (h *MonitorHandler) GetMonitorById(c echo.Context) error {
	ctx := c.Request().Context()
	log := LoggerFromContext(c)
	userId := GetUserID(c)

	monitorId := c.Param("monitorId")
	id, err := uuid.Parse(monitorId)
	if err != nil {
		log.Error().Err(err).Msg("invalid monitor id")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid monitor id")
	}

	monitor, err := h.monitorService.GetMonitorById(ctx, userId, id)
	if err != nil {
		log.Error().Err(err).Str("id", id.String()).Msg("failed to get monitor with id")
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get monitor")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    map[string]any{"monitor": monitor},
	})
}