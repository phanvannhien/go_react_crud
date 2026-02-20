package test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/internal/handler"
	"github.com/nhienphan/full_rest_app/internal/middleware"
)

// TestUserEndpoints_Unauthorized verifies that all user endpoints require auth.
func TestUserEndpoints_Unauthorized(t *testing.T) {
	e := echo.New()

	userHandler := handler.NewUserHandler(nil)
	api := e.Group("/api", middleware.JWTAuth("test-secret"))
	api.GET("/users", userHandler.ListUsers)
	api.GET("/users/:id", userHandler.GetUser)
	api.PUT("/users/:id", userHandler.UpdateUser)
	api.DELETE("/users/:id", userHandler.DeleteUser)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"list users", http.MethodGet, "/api/users"},
		{"get user", http.MethodGet, "/api/users/some-uuid"},
		{"update user", http.MethodPut, "/api/users/some-uuid"},
		{"delete user", http.MethodDelete, "/api/users/some-uuid"},
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

// TestUpdateUser_ValidationErrors tests DTO validation for user updates.
func TestUpdateUser_ValidationErrors(t *testing.T) {
	e := echo.New()
	userHandler := handler.NewUserHandler(nil)
	e.PUT("/api/users/:id", userHandler.UpdateUser) // No auth for validation test

	tests := []struct {
		name     string
		body     map[string]interface{}
		wantCode int
	}{
		{
			name:     "invalid email format",
			body:     map[string]interface{}{"email": "not-an-email"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid role value",
			body:     map[string]interface{}{"role": "superadmin"},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(e, http.MethodPut, "/api/users/test-id", tt.body, "")
			assertStatus(t, rec, tt.wantCode)

			result := parseResponse(t, rec)
			assertHasError(t, result)
		})
	}
}
