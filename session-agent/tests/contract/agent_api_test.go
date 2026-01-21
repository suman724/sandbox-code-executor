package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"session-agent/internal/api"
	"session-agent/internal/api/middleware"
)

func TestAgentRoutes(t *testing.T) {
	stepsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if payload["sessionId"] == "" || payload["stepId"] == "" || payload["code"] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"completed","stdout":"","stderr":""}`))
	})

	router := api.NewRouter(api.RouterDeps{StepsHandler: stepsHandler})

	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	body, err := json.Marshal(map[string]string{
		"sessionId": "session-1",
		"stepId":    "step-1",
		"code":      "print(1)",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	req = httptest.NewRequest(http.MethodPost, "/v1/steps", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestAgentAuthMiddleware(t *testing.T) {
	stepsHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router := api.NewRouter(api.RouterDeps{
		StepsHandler:  stepsHandler,
		AuthMiddleware: middleware.AuthTokenMiddleware("token"),
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/steps", bytes.NewReader([]byte(`{"sessionId":"s","stepId":"t","code":"print(1)"}`)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/v1/steps", bytes.NewReader([]byte(`{"sessionId":"s","stepId":"t","code":"print(1)"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-Token", "token")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
