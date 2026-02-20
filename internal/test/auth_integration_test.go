package test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nhienphan/full_rest_app/db/sqlc"
	"github.com/nhienphan/full_rest_app/internal/handler"
	"github.com/nhienphan/full_rest_app/internal/routes"
	"github.com/nhienphan/full_rest_app/internal/service"
)

// NOTE: These integration tests require a running PostgreSQL instance.
// Set DATABASE_URL environment variable before running.
// Example: DATABASE_URL=postgres://user:pass@localhost:5432/test_db?sslmode=disable

// setupAuthTestServer creates an Echo instance wired for auth tests.
// In a real setup, this would connect to a test database.
func setupAuthTestServer(queries *sqlc.Queries, jwtSecret string) *echo.Echo {
	authService := service.NewAuthService(queries, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	e := newTestEcho()
	routes.RegisterAuthRoutes(e, authHandler)
	return e
}

func TestRegister_ValidationErrors(t *testing.T) {
	// Test without database - validates handler DTO validation
	e := echo.New()
	authHandler := handler.NewAuthHandler(nil) // nil service, we only test validation
	e.POST("/api/auth/register", authHandler.Register)

	tests := []struct {
		name     string
		body     map[string]interface{}
		wantCode int
	}{
		{
			name:     "missing email",
			body:     map[string]interface{}{"password": "12345678"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid email format",
			body:     map[string]interface{}{"email": "not-an-email", "password": "12345678"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "password too short",
			body:     map[string]interface{}{"email": "test@example.com", "password": "123"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty body",
			body:     map[string]interface{}{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid role",
			body:     map[string]interface{}{"email": "test@example.com", "password": "12345678", "role": "superadmin"},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(e, http.MethodPost, "/api/auth/register", tt.body, "")
			assertStatus(t, rec, tt.wantCode)

			result := parseResponse(t, rec)
			assertHasError(t, result)
		})
	}
}

func TestLogin_ValidationErrors(t *testing.T) {
	e := echo.New()
	authHandler := handler.NewAuthHandler(nil)
	e.POST("/api/auth/login", authHandler.Login)

	tests := []struct {
		name     string
		body     map[string]interface{}
		wantCode int
	}{
		{
			name:     "missing email",
			body:     map[string]interface{}{"password": "12345678"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid email",
			body:     map[string]interface{}{"email": "bad", "password": "12345678"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing password",
			body:     map[string]interface{}{"email": "test@example.com"},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(e, http.MethodPost, "/api/auth/login", tt.body, "")
			assertStatus(t, rec, tt.wantCode)

			result := parseResponse(t, rec)
			assertHasError(t, result)
		})
	}
}
