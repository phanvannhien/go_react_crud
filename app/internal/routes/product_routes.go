package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/handler"
)

// RegisterProductRoutes registers product routes (auth required).
func RegisterProductRoutes(g *echo.Group, h *handler.ProductHandler) {
	products := g.Group("/products")
	products.GET("", h.ListProducts)
	products.GET("/:id", h.GetProduct)
	products.POST("", h.CreateProduct)
	products.PUT("/:id", h.UpdateProduct)
	products.DELETE("/:id", h.DeleteProduct)
}
