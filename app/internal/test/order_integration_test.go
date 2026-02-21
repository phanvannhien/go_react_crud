package test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/handler"
	"github.com/nhienphan/full_rest_app/internal/middleware"
)

// TestOrderEndpoints_Unauthorized verifies all order endpoints require auth.
func TestOrderEndpoints_Unauthorized(t *testing.T) {
	e := echo.New()

	orderHandler := handler.NewOrderHandler(nil)
	api := e.Group("/api", middleware.JWTAuth("test-secret"))
	api.GET("/orders", orderHandler.ListOrders)
	api.GET("/orders/:id", orderHandler.GetOrder)
	api.POST("/orders", orderHandler.CreateOrder)
	api.PUT("/orders/:id", orderHandler.UpdateOrderStatus)
	api.DELETE("/orders/:id", orderHandler.DeleteOrder)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"list orders", http.MethodGet, "/api/orders"},
		{"get order", http.MethodGet, "/api/orders/some-uuid"},
		{"create order", http.MethodPost, "/api/orders"},
		{"update order", http.MethodPut, "/api/orders/some-uuid"},
		{"delete order", http.MethodDelete, "/api/orders/some-uuid"},
	}

	for _, tt := range tests {
		t.Run(tt.name+" without token", func(t *testing.T) {
			rec := doRequest(e, tt.method, tt.path, nil, "")
			assertStatus(t, rec, http.StatusUnauthorized)
		})

		t.Run(tt.name+" with invalid token", func(t *testing.T) {
			rec := doRequest(e, tt.method, tt.path, nil, "invalid-token")
			assertStatus(t, rec, http.StatusUnauthorized)
		})
	}
}

// fakeAuthMiddleware injects a fake user_id into context for validation testing.
func fakeAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("user_id", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
		return next(c)
	}
}

// TestCreateOrder_ValidationErrors tests DTO validation for order creation.
func TestCreateOrder_ValidationErrors(t *testing.T) {
	e := echo.New()
	orderHandler := handler.NewOrderHandler(nil)
	g := e.Group("", fakeAuthMiddleware)
	g.POST("/api/orders", orderHandler.CreateOrder)

	tests := []struct {
		name     string
		body     map[string]interface{}
		wantCode int
	}{
		{
			name:     "empty body",
			body:     map[string]interface{}{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty items array",
			body:     map[string]interface{}{"items": []interface{}{}},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "missing product_id in item",
			body: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"quantity": 1, "price": 10.0},
				},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "invalid product_id format",
			body: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"product_id": "not-a-uuid", "quantity": 1, "price": 10.0},
				},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "zero quantity",
			body: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"product_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "quantity": 0, "price": 10.0},
				},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "negative price",
			body: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"product_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "quantity": 1, "price": -5.0},
				},
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(e, http.MethodPost, "/api/orders", tt.body, "")
			assertStatus(t, rec, tt.wantCode)

			result := parseResponse(t, rec)
			assertHasError(t, result)
		})
	}
}

// TestUpdateOrderStatus_ValidationErrors tests DTO validation for status update.
func TestUpdateOrderStatus_ValidationErrors(t *testing.T) {
	e := echo.New()
	orderHandler := handler.NewOrderHandler(nil)
	g := e.Group("", fakeAuthMiddleware)
	g.PUT("/api/orders/:id", orderHandler.UpdateOrderStatus)

	tests := []struct {
		name     string
		body     map[string]interface{}
		wantCode int
	}{
		{
			name:     "empty body",
			body:     map[string]interface{}{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid status value",
			body:     map[string]interface{}{"status": "unknown"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty status",
			body:     map[string]interface{}{"status": ""},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(e, http.MethodPut, "/api/orders/some-uuid", tt.body, "")
			assertStatus(t, rec, tt.wantCode)

			result := parseResponse(t, rec)
			assertHasError(t, result)
		})
	}
}
