package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"session-agent/internal/api/middleware"
)

func TestAuthBypassAllowsRequest(t *testing.T) {
	handler := middleware.AuthBypassMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/steps", nil)
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestAuthTokenMiddlewareRejectsMissingToken(t *testing.T) {
	handler := middleware.AuthTokenMiddleware("token")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/steps", nil)
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}
