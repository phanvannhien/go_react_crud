package handler

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/response"
	"github.com/nhienphan/full_rest_app/internal/service"
)

// AuthHandler handles auth HTTP endpoints.
type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// RegisterRequest is the DTO for registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=128"`
	Role     string `json:"role" validate:"omitempty,oneof=user admin"`
}

// LoginRequest is the DTO for login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Register handles POST /api/auth/register.
func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
	}

	// Trim strings
	req.Email = strings.TrimSpace(req.Email)
	req.Email = strings.ToLower(req.Email)

	// Validate
	if err := h.validate.Struct(req); err != nil {
		errors := formatValidationErrors(err)
		return response.ValidationError(c, errors)
	}

	// Default role
	if req.Role == "" {
		req.Role = "user"
	}

	result, err := h.authService.Register(c.Request().Context(), req.Email, req.Password, req.Role)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return response.Conflict(c, "Email already exists")
		}
		return response.InternalError(c)
	}

	return response.Created(c, result)
}

// Login handles POST /api/auth/login.
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad_request", "Invalid request body")
	}

	// Trim strings
	req.Email = strings.TrimSpace(req.Email)
	req.Email = strings.ToLower(req.Email)

	// Validate
	if err := h.validate.Struct(req); err != nil {
		errors := formatValidationErrors(err)
		return response.ValidationError(c, errors)
	}

	result, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", "Invalid email or password")
	}

	return response.Success(c, result)
}

// formatValidationErrors converts validator errors to a map.
func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				errors[field] = field + " must be at least " + e.Param() + " characters"
			case "max":
				errors[field] = field + " must be at most " + e.Param() + " characters"
			case "oneof":
				errors[field] = field + " must be one of: " + e.Param()
			case "uuid":
				errors[field] = field + " must be a valid UUID"
			default:
				errors[field] = field + " is invalid"
			}
		}
	}
	return errors
}
