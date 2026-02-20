package test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/handler"
	"github.com/nhienphan/full_rest_app/internal/middleware"
)

// TestProductEndpoints_Unauthorized verifies all product endpoints require auth.
func TestProductEndpoints_Unauthorized(t *testing.T) {
	e := echo.New()

	productHandler := handler.NewProductHandler(nil)
	api := e.Group("/api", middleware.JWTAuth("test-secret"))
	api.GET("/products", productHandler.ListProducts)
	api.GET("/products/:id", productHandler.GetProduct)
	api.POST("/products", productHandler.CreateProduct)
	api.PUT("/products/:id", productHandler.UpdateProduct)
	api.DELETE("/products/:id", productHandler.DeleteProduct)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"list products", http.MethodGet, "/api/products"},
		{"get product", http.MethodGet, "/api/products/some-uuid"},
		{"create product", http.MethodPost, "/api/products"},
		{"update product", http.MethodPut, "/api/products/some-uuid"},
		{"delete product", http.MethodDelete, "/api/products/some-uuid"},
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

// TestCreateProduct_ValidationErrors tests DTO validation for product creation.
func TestCreateProduct_ValidationErrors(t *testing.T) {
	e := echo.New()
	productHandler := handler.NewProductHandler(nil)
	e.POST("/api/products", productHandler.CreateProduct)

	tests := []struct {
		name     string
		body     map[string]interface{}
		wantCode int
	}{
		{
			name:     "missing name",
			body:     map[string]interface{}{"price": 10, "category_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing price",
			body:     map[string]interface{}{"name": "Test", "category_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing category_id",
			body:     map[string]interface{}{"name": "Test", "price": 10},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid category_id format",
			body:     map[string]interface{}{"name": "Test", "price": 10, "category_id": "not-a-uuid"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "negative price",
			body:     map[string]interface{}{"name": "Test", "price": -5, "category_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty body",
			body:     map[string]interface{}{},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(e, http.MethodPost, "/api/products", tt.body, "")
			assertStatus(t, rec, tt.wantCode)

			result := parseResponse(t, rec)
			assertHasError(t, result)
		})
	}
}

// TestListProducts_QueryParamValidation tests query param validation.
func TestListProducts_QueryParamValidation(t *testing.T) {
	e := echo.New()
	productHandler := handler.NewProductHandler(nil)
	e.GET("/api/products", productHandler.ListProducts)

	tests := []struct {
		name     string
		query    string
		wantCode int
	}{
		{
			name:     "invalid is_active",
			query:    "?is_active=notabool",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid min_stock",
			query:    "?min_stock=abc",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid min_price",
			query:    "?min_price=abc",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid max_price",
			query:    "?max_price=abc",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(e, http.MethodGet, "/api/products"+tt.query, nil, "")
			assertStatus(t, rec, tt.wantCode)

			result := parseResponse(t, rec)
			assertHasError(t, result)
		})
	}
}
