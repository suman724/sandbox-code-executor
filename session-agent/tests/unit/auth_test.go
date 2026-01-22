package unit

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"session-agent/internal/api/handlers"
	"session-agent/internal/api/middleware"
	"session-agent/internal/runtime"
	"shared/sessionagent"
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

func TestRequireTokenMiddlewareRejectsMissingToken(t *testing.T) {
	handler := middleware.RequireTokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/steps", nil)
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestPerSessionTokenEnforcement(t *testing.T) {
	runner := runtime.NewRunner()
	_, err := runner.RegisterSession(sessionagent.SessionRegisterRequest{
		SessionID: "session-1",
		Runtime:   "python",
		Token:     "token-1",
	})
	if err != nil {
		t.Fatalf("register session: %v", err)
	}
	defer runner.RemoveSession("session-1")
	handler := handlers.StepHandler{Runner: runner, RequireToken: true}
	body := []byte(`{"sessionId":"session-1","stepId":"step-1","code":"print(1)"}`)

	req := httptest.NewRequest(http.MethodPost, "/v1/steps", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-Token", "wrong")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/v1/steps", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-Token", "token-1")
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
