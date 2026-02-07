package middleware

import (
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/labstack/echo/v4"

	"github.com/pulkitbhatt/ikiru/internal/server"
	"github.com/pulkitbhatt/ikiru/internal/service"
)

type contextKey string

const (
	ContextUserID    contextKey = "user_id"
	ContextIDPUserID contextKey = "idp_user_id"
)

type AuthMiddleware struct {
	server      *server.Server
	authService *service.AuthService
}

func NewAuthMiddleware(s *server.Server, a *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		server:      s,
		authService: a,
	}
}

func (auth *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return (func(c echo.Context) error {
		start := time.Now()
		claims, ok := clerk.SessionClaimsFromContext(c.Request().Context())
		if !ok {
			auth.server.Logger.Error().
				Str("function", "RequireAuth").
				Dur("duration", time.Since(start)).
				Msg("failed to get session from context")
			return echo.ErrUnauthorized
		}

		userID, err := auth.authService.EnsureUser(c.Request().Context(), claims.Subject, "")
		if err != nil {
			auth.server.Logger.Error().
				Err(err).
				Str("function", "RequireAuth").
				Dur("duration", time.Since(start)).
				Str("idp_user_id", claims.Subject).
				Msg("failed to ensure user")
			return echo.ErrInternalServerError
		}

		auth.server.Logger.Info().Msgf("userid in auth: %v", userID)

		c.Set(string(ContextIDPUserID), claims.Subject)
		c.Set(string(ContextUserID), userID)
		return next(c)
	})
}
