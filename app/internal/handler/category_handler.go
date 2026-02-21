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

// CategoryHandler handles category HTTP endpoints.
type CategoryHandler struct {
	categoryService *service.CategoryService
	validate        *validator.Validate
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validate:        validator.New(),
	}
}

// CreateCategoryRequest is the DTO for creating a category.
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

// UpdateCategoryRequest is the DTO for updating a category.
type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

// ListCategories handles GET /api/categories.
func (h *CategoryHandler) ListCategories(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	result, err := h.categoryService.ListCategories(c.Request().Context(), page, limit)
	if err != nil {
		return response.InternalError(c)
	}

	return response.Paginated(c, result.Categories, page, limit, result.Total)
}

// GetCategory handles GET /api/categories/:id.
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "Category ID is required"})
	}

	cat, err := h.categoryService.GetCategoryByID(c.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			return response.NotFound(c, "Category")
		}
		return response.InternalError(c)
	}

	return response.Success(c, cat)
}

// CreateCategory handles POST /api/categories.
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var req CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		errors := formatValidationErrors(err)
		return response.ValidationError(c, errors)
	}

	cat, err := h.categoryService.CreateCategory(c.Request().Context(), strings.TrimSpace(req.Name))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return response.Error(c, http.StatusConflict, "conflict", "Category name already exists")
		}
		return response.InternalError(c)
	}

	return response.Created(c, cat)
}

// UpdateCategory handles PUT /api/categories/:id.
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "Category ID is required"})
	}

	var req UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		errors := formatValidationErrors(err)
		return response.ValidationError(c, errors)
	}

	cat, err := h.categoryService.UpdateCategory(c.Request().Context(), id, strings.TrimSpace(req.Name))
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			return response.NotFound(c, "Category")
		}
		return response.InternalError(c)
	}

	return response.Success(c, cat)
}

// DeleteCategory handles DELETE /api/categories/:id.
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "Category ID is required"})
	}

	err := h.categoryService.DeleteCategory(c.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			return response.NotFound(c, "Category")
		}
		return response.InternalError(c)
	}

	return response.Success(c, map[string]string{"message": "Category deleted successfully"})
}
