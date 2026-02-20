package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIResponse is the standard response envelope.
type APIResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

// PaginatedResponse wraps list data with offset pagination metadata.
type PaginatedResponse struct {
	Data  interface{}    `json:"data"`
	Meta  PaginationMeta `json:"meta"`
	Error interface{}    `json:"error"`
}

// PaginationMeta holds offset pagination metadata.
type PaginationMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// CursorResponse wraps list data with cursor pagination metadata.
type CursorResponse struct {
	Data       interface{} `json:"data"`
	NextCursor string      `json:"next_cursor"`
	Error      interface{} `json:"error"`
}

// Success returns a 200 JSON response with data.
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, APIResponse{
		Data:  data,
		Error: nil,
	})
}

// Created returns a 201 JSON response with data.
func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, APIResponse{
		Data:  data,
		Error: nil,
	})
}

// Paginated returns a 200 JSON response with data and offset pagination meta.
func Paginated(c echo.Context, data interface{}, page, limit int, total int64) error {
	return c.JSON(http.StatusOK, PaginatedResponse{
		Data: data,
		Meta: PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
		Error: nil,
	})
}

// Cursor returns a 200 JSON response with data and next cursor.
func Cursor(c echo.Context, data interface{}, nextCursor string) error {
	return c.JSON(http.StatusOK, CursorResponse{
		Data:       data,
		NextCursor: nextCursor,
		Error:      nil,
	})
}

// ErrorResponse is the standard error payload.
type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error returns an error JSON response.
func Error(c echo.Context, status int, code, message string) error {
	return c.JSON(status, APIResponse{
		Data: nil,
		Error: ErrorResponse{
			Code:    code,
			Message: message,
		},
	})
}

// ValidationError returns a 400 validation error JSON response.
func ValidationError(c echo.Context, details interface{}) error {
	return c.JSON(http.StatusBadRequest, APIResponse{
		Data: nil,
		Error: ErrorResponse{
			Code:    "validation_error",
			Message: "Validation failed",
			Details: details,
		},
	})
}

// InternalError returns a 500 generic internal error.
func InternalError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, APIResponse{
		Data: nil,
		Error: ErrorResponse{
			Code:    "internal_error",
			Message: "An internal error occurred",
		},
	})
}

// NotFound returns a 404 not found error.
func NotFound(c echo.Context, resource string) error {
	return Error(c, http.StatusNotFound, "not_found", resource+" not found")
}

// Unauthorized returns a 401 error.
func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, "unauthorized", message)
}

// Conflict returns a 409 error.
func Conflict(c echo.Context, message string) error {
	return Error(c, http.StatusConflict, "conflict", message)
}
