package middleware

import (
	"github.com/pulkitbhatt/ikiru/internal/handler"
	"github.com/pulkitbhatt/ikiru/internal/server"
)

type Middlewares struct {
	Auth            AuthMiddleware
	ContextEnhancer ContextEnhancer
}

func NewMiddlewares(s *server.Server, h *handler.Handlers) *Middlewares {
	return &Middlewares{
		Auth:            *NewAuthMiddleware(s, &h.Services.Auth),
		ContextEnhancer: *NewContextEnhancer(s),
	}
}
