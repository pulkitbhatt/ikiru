package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/pulkitbhatt/ikiru/internal/server"
)

const (
	UserIDKey = "user_id"
	LoggerKey = "logger"
)

type ContextEnhancer struct {
	server *server.Server
}

func NewContextEnhancer(s *server.Server) *ContextEnhancer {
	return &ContextEnhancer{
		server: s,
	}
}

func (ce *ContextEnhancer) EnhanceContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := GetRequestID(c)

			contextLogger := ce.server.Logger.With().
				Str(RequestIDKey, id).
				Str("method", c.Request().Method).
				Str("path", c.Path()).
				Logger()

			if userId := ce.extractUserID(c); userId != "" {
				contextLogger = contextLogger.With().Str(UserIDKey, userId).Logger()
			}

			ctx := contextLogger.WithContext(c.Request().Context())
			c.Set(LoggerKey, &contextLogger)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func (ce *ContextEnhancer) extractUserID(c echo.Context) string {
	if userID, ok := c.Get(UserIDKey).(string); ok && userID != "" {
		return userID
	}
	return ""
}
