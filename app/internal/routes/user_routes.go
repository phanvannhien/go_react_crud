package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/handler"
)

// RegisterUserRoutes registers user routes (auth required).
func RegisterUserRoutes(g *echo.Group, h *handler.UserHandler) {
	users := g.Group("/users")
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetUser)
	users.PUT("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
}
