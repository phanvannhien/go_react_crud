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

// UserHandler handles user HTTP endpoints.
type UserHandler struct {
	userService *service.UserService
	validate    *validator.Validate
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// UpdateUserRequest is the DTO for updating a user.
type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email,max=255"`
	Role     string `json:"role" validate:"omitempty,oneof=user admin"`
	IsActive *bool  `json:"is_active" validate:"omitempty"`
}

// ListUsers handles GET /api/users.
func (h *UserHandler) ListUsers(c echo.Context) error {
	// Parse query params
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	sort := strings.TrimSpace(c.QueryParam("sort"))
	order := strings.TrimSpace(c.QueryParam("order"))

	var search *string
	if s := strings.TrimSpace(c.QueryParam("search")); s != "" {
		search = &s
	}

	var role *string
	if r := strings.TrimSpace(c.QueryParam("role")); r != "" {
		role = &r
	}

	var isActive *bool
	if ia := strings.TrimSpace(c.QueryParam("is_active")); ia != "" {
		val, err := strconv.ParseBool(ia)
		if err != nil {
			return response.ValidationError(c, map[string]string{"is_active": "must be true or false"})
		}
		isActive = &val
	}

	params := service.UserListParams{
		Search:    search,
		Role:      role,
		IsActive:  isActive,
		SortField: sort,
		SortOrder: order,
		Page:      page,
		Limit:     limit,
	}

	result, err := h.userService.ListUsers(c.Request().Context(), params)
	if err != nil {
		return response.InternalError(c)
	}

	return response.Paginated(c, result.Users, params.Page, params.Limit, result.Total)
}

// GetUser handles GET /api/users/:id.
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "User ID is required"})
	}

	user, err := h.userService.GetUser(c.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			return response.NotFound(c, "User")
		}
		return response.InternalError(c)
	}

	return response.Success(c, user)
}

// UpdateUser handles PUT /api/users/:id.
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "User ID is required"})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
	}

	// Trim
	req.Email = strings.TrimSpace(req.Email)
	if req.Email != "" {
		req.Email = strings.ToLower(req.Email)
	}

	// Validate
	if err := h.validate.Struct(req); err != nil {
		errors := formatValidationErrors(err)
		return response.ValidationError(c, errors)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), id, req.Email, req.Role, isActive)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			return response.NotFound(c, "User")
		}
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return response.Conflict(c, "Email already exists")
		}
		return response.InternalError(c)
	}

	return response.Success(c, user)
}

// DeleteUser handles DELETE /api/users/:id.
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.ValidationError(c, map[string]string{"id": "User ID is required"})
	}

	err := h.userService.DeleteUser(c.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(c, "User")
		}
		return response.InternalError(c)
	}

	return response.Success(c, map[string]string{"message": "User deleted successfully"})
}
