package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/middleware"
	"github.com/nhienphan/full_rest_app/internal/response"
	"github.com/nhienphan/full_rest_app/internal/service"
)

// OrderHandler handles order HTTP endpoints.
type OrderHandler struct {
	orderService *service.OrderService
	validate     *validator.Validate
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		validate:     validator.New(),
	}
}

// CreateOrderItemRequest is the DTO for an order item.
type CreateOrderItemRequest struct {
	ProductID string  `json:"product_id" validate:"required,uuid"`
	Quantity  int32   `json:"quantity" validate:"required,gt=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

// CreateOrderRequest is the DTO for creating an order.
type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// UpdateOrderStatusRequest is the DTO for updating order status.
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending paid cancelled shipped"`
}

// CreateOrder handles POST /api/orders.
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
	}

	// Validate
	if err := h.validate.Struct(req); err != nil {
		errors := formatValidationErrors(err)
		return response.ValidationError(c, errors)
	}

	// Map to service params
	items := make([]service.OrderCreateItemParams, len(req.Items))
	for i, item := range req.Items {
		items[i] = service.OrderCreateItemParams{
			ProductID: strings.TrimSpace(item.ProductID),
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	order, err := h.orderService.CreateOrder(c.Request().Context(), service.OrderCreateParams{
		UserID: userID,
		Items:  items,
	})
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			return response.ValidationError(c, map[string]string{"error": err.Error()})
		}
		return response.InternalError(c)
	}

	return response.Created(c, order)
}

// GetOrder handles GET /api/orders/:id.
func (h *OrderHandler) GetOrder(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return response.Unauthorized(c, "User not authenticated")
	}

	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "Order ID is required"})
	}

	order, err := h.orderService.GetOrderByID(c.Request().Context(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			return response.NotFound(c, "Order")
		}
		if strings.Contains(err.Error(), "forbidden") {
			return response.Error(c, http.StatusForbidden, "forbidden", err.Error())
		}
		return response.InternalError(c)
	}

	return response.Success(c, order)
}

// UpdateOrderStatus handles PUT /api/orders/:id.
func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return response.Unauthorized(c, "User not authenticated")
	}

	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "Order ID is required"})
	}

	var req UpdateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
	}

	// Validate
	if err := h.validate.Struct(req); err != nil {
		errors := formatValidationErrors(err)
		return response.ValidationError(c, errors)
	}

	order, err := h.orderService.UpdateOrderStatus(c.Request().Context(), id, userID, req.Status)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			return response.NotFound(c, "Order")
		}
		if strings.Contains(err.Error(), "forbidden") {
			return response.Error(c, http.StatusForbidden, "forbidden", err.Error())
		}
		if strings.Contains(err.Error(), "cannot update") {
			return response.Error(c, http.StatusConflict, "conflict", err.Error())
		}
		return response.InternalError(c)
	}

	return response.Success(c, order)
}

// DeleteOrder handles DELETE /api/orders/:id.
func (h *OrderHandler) DeleteOrder(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return response.Unauthorized(c, "User not authenticated")
	}

	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "Order ID is required"})
	}

	err := h.orderService.DeleteOrder(c.Request().Context(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			return response.NotFound(c, "Order")
		}
		if strings.Contains(err.Error(), "forbidden") {
			return response.Error(c, http.StatusForbidden, "forbidden", err.Error())
		}
		if strings.Contains(err.Error(), "cannot delete") {
			return response.Error(c, http.StatusConflict, "conflict", err.Error())
		}
		return response.InternalError(c)
	}

	return response.Success(c, map[string]string{"message": "Order deleted successfully"})
}

// ListOrders handles GET /api/orders.
func (h *OrderHandler) ListOrders(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return response.Unauthorized(c, "User not authenticated")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	result, err := h.orderService.ListOrders(c.Request().Context(), userID, page, limit)
	if err != nil {
		return response.InternalError(c)
	}

	return response.Paginated(c, result.Orders, page, limit, result.Total)
}
