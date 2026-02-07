package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/pulkitbhatt/ikiru/internal/util"
)

const (
	RequestIDKey    = "request_id"
	RequestIDHeader = "X-Request-ID"
)

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rqstId := util.GenerateUUID()
			c.Set(RequestIDKey, rqstId)
			c.Response().Header().Set(RequestIDHeader, rqstId)
			return next(c)
		}
	}
}

func GetRequestID(c echo.Context) string {
	if id, ok := c.Get(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
