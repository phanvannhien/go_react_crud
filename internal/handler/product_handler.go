package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/response"
	"github.com/nhienphan/full_rest_app/internal/service"
)

// ProductHandler handles product HTTP endpoints.
type ProductHandler struct {
	productService *service.ProductService
	validate       *validator.Validate
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		validate:       validator.New(),
	}
}

// CreateProductRequest is the DTO for creating a product.
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Description string  `json:"description" validate:"max=5000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	Stock       int32   `json:"stock" validate:"gte=0"`
	IsActive    *bool   `json:"is_active"`
}

// UpdateProductRequest is the DTO for updating a product.
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"omitempty,max=255"`
	Description string  `json:"description" validate:"max=5000"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	CategoryID  string  `json:"category_id" validate:"omitempty,uuid"`
	Stock       *int32  `json:"stock" validate:"omitempty,gte=0"`
	IsActive    *bool   `json:"is_active"`
}

// ListProducts handles GET /api/products (cursor pagination).
func (h *ProductHandler) ListProducts(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	var search *string
	if s := strings.TrimSpace(c.QueryParam("search")); s != "" {
		search = &s
	}

	var categoryID *string
	if cat := strings.TrimSpace(c.QueryParam("category_id")); cat != "" {
		categoryID = &cat
	}

	var isActive *bool
	if ia := strings.TrimSpace(c.QueryParam("is_active")); ia != "" {
		val, err := strconv.ParseBool(ia)
		if err != nil {
			return response.ValidationError(c, map[string]string{"is_active": "must be true or false"})
		}
		isActive = &val
	}

	var minStock *int32
	if ms := strings.TrimSpace(c.QueryParam("min_stock")); ms != "" {
		val, err := strconv.ParseInt(ms, 10, 32)
		if err != nil {
			return response.ValidationError(c, map[string]string{"min_stock": "must be an integer"})
		}
		v := int32(val)
		minStock = &v
	}

	var minPrice *float64
	if mp := strings.TrimSpace(c.QueryParam("min_price")); mp != "" {
		val, err := strconv.ParseFloat(mp, 64)
		if err != nil {
			return response.ValidationError(c, map[string]string{"min_price": "must be a number"})
		}
		minPrice = &val
	}

	var maxPrice *float64
	if mp := strings.TrimSpace(c.QueryParam("max_price")); mp != "" {
		val, err := strconv.ParseFloat(mp, 64)
		if err != nil {
			return response.ValidationError(c, map[string]string{"max_price": "must be a number"})
		}
		maxPrice = &val
	}

	var cursor *string
	if cur := strings.TrimSpace(c.QueryParam("cursor")); cur != "" {
		cursor = &cur
	}

	params := service.ProductListParams{
		Search:     search,
		CategoryID: categoryID,
		IsActive:   isActive,
		MinStock:   minStock,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Cursor:     cursor,
		Limit:      limit,
	}

	result, err := h.productService.ListProducts(c.Request().Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			return response.ValidationError(c, map[string]string{"error": err.Error()})
		}
		return response.InternalError(c)
	}

	return response.Cursor(c, result.Products, result.NextCursor)
}

// GetProduct handles GET /api/products/:id.
func (h *ProductHandler) GetProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "Product ID is required"})
	}

	product, err := h.productService.GetProduct(c.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			return response.NotFound(c, "Product")
		}
		return response.InternalError(c)
	}

	return response.Success(c, product)
}

// CreateProduct handles POST /api/products.
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
	}

	// Trim
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.CategoryID = strings.TrimSpace(req.CategoryID)

	// Validate
	if err := h.validate.Struct(req); err != nil {
		errors := formatValidationErrors(err)
		return response.ValidationError(c, errors)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	product, err := h.productService.CreateProduct(c.Request().Context(), service.ProductCreateParams{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		Stock:       req.Stock,
		IsActive:    isActive,
	})
	if err != nil {
		if strings.Contains(err.Error(), "invalid category") {
			return response.ValidationError(c, map[string]string{"category_id": "Invalid category ID"})
		}
		return response.InternalError(c)
	}

	return response.Created(c, product)
}

// UpdateProduct handles PUT /api/products/:id.
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "Product ID is required"})
	}

	var req UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
	}

	// Trim
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.CategoryID = strings.TrimSpace(req.CategoryID)

	// Validate
	if err := h.validate.Struct(req); err != nil {
		errors := formatValidationErrors(err)
		return response.ValidationError(c, errors)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	stock := int32(0)
	if req.Stock != nil {
		stock = *req.Stock
	}

	product, err := h.productService.UpdateProduct(c.Request().Context(), service.ProductUpdateParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		Stock:       stock,
		IsActive:    isActive,
	})
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			return response.NotFound(c, "Product")
		}
		return response.InternalError(c)
	}

	return response.Success(c, product)
}

// DeleteProduct handles DELETE /api/products/:id.
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "Product ID is required"})
	}

	err := h.productService.DeleteProduct(c.Request().Context(), id)
	if err != nil {
		return response.InternalError(c)
	}

	return response.Success(c, map[string]string{"message": "Product deleted successfully"})
}
