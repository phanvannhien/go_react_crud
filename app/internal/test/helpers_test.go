package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

// Helper functions for integration tests

func newTestEcho() *echo.Echo {
	e := echo.New()
	return e
}

func doRequest(e *echo.Echo, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if token != "" {
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func parseResponse(t *testing.T, rec *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v, body: %s", err, rec.Body.String())
	}
	return result
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rec.Code != expected {
		t.Errorf("Expected status %d, got %d. Body: %s", expected, rec.Code, rec.Body.String())
	}
}

func assertNoError(t *testing.T, result map[string]interface{}) {
	t.Helper()
	if result["error"] != nil {
		t.Errorf("Expected no error, got: %v", result["error"])
	}
}

func assertHasError(t *testing.T, result map[string]interface{}) {
	t.Helper()
	if result["error"] == nil {
		t.Error("Expected error, got nil")
	}
}
