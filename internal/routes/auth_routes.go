package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/handler"
)

// RegisterAuthRoutes registers auth routes (no auth required).
func RegisterAuthRoutes(e *echo.Echo, h *handler.AuthHandler) {
	auth := e.Group("/api/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
}
