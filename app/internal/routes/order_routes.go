package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/handler"
)

// RegisterOrderRoutes registers order routes (auth required).
func RegisterOrderRoutes(g *echo.Group, h *handler.OrderHandler) {
	orders := g.Group("/orders")
	orders.GET("", h.ListOrders)
	orders.GET("/:id", h.GetOrder)
	orders.POST("", h.CreateOrder)
	orders.PUT("/:id", h.UpdateOrderStatus)
	orders.DELETE("/:id", h.DeleteOrder)
}
