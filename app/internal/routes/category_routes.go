package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/handler"
)

// RegisterCategoryRoutes registers category routes (auth required).
func RegisterCategoryRoutes(g *echo.Group, h *handler.CategoryHandler) {
	categories := g.Group("/categories")
	categories.GET("", h.ListCategories)
	categories.GET("/:id", h.GetCategory)
	categories.POST("", h.CreateCategory)
	categories.PUT("/:id", h.UpdateCategory)
	categories.DELETE("/:id", h.DeleteCategory)
}
