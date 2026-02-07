package router

import (
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/labstack/echo/v4"

	"github.com/pulkitbhatt/ikiru/internal/handler"
	"github.com/pulkitbhatt/ikiru/internal/middleware"
	"github.com/pulkitbhatt/ikiru/internal/server"
)

func NewRouter(s *server.Server, h *handler.Handlers) *echo.Echo {
	middlewares := middleware.NewMiddlewares(s, h)
	router := echo.New()
	api := router.Group("/v1",
		middleware.RequestID(),
		echo.WrapMiddleware(clerkhttp.WithHeaderAuthorization()),
	)

	public := api.Group("")

	pvt := api.Group("",
		middlewares.Auth.RequireAuth,
		middlewares.ContextEnhancer.EnhanceContext(),
	)

	public.GET("/health", h.Health.CheckHealth)
	pvt.POST("/monitor", h.Monitor.CreateMonitor)
	return router
}
