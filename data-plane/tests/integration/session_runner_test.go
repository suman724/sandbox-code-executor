package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"data-plane/internal/execution"
	"data-plane/internal/runtime"
)

func TestSessionRunnerLocal(t *testing.T) {
	t.Setenv("AUTHZ_BYPASS", "true")

	sessionRuntime := runtime.NewLocalSessionRuntime()
	sessionRegistry := runtime.NewInMemorySessionRegistry()
	sessionHandler := runtime.SessionHandler{Runtime: sessionRuntime, Registry: sessionRegistry}

	runHandler := runtime.RunHandler{
		Runner: execution.Runner{
			Registry: runtime.DefaultRegistry(),
			Deps:     runtime.DependencyPolicy{},
		},
	}

	handler := runtime.RouterWithDependencies(runtime.Dependencies{
		RunHandler:     runHandler,
		SessionHandler: sessionHandler,
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	payload := map[string]any{
		"sessionId":    "session-1",
		"policyId":     "policy-1",
		"workspaceRef": "session-1",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	resp, err := http.Post(server.URL+"/sessions", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, resp.StatusCode)
	}
	_ = resp.Body.Close()

	stepBody, err := json.Marshal(map[string]string{"command": "echo out; echo err 1>&2"})
	if err != nil {
		t.Fatalf("marshal step payload: %v", err)
	}
	resp, err = http.Post(server.URL+"/sessions/session-1/steps", "application/json", bytes.NewReader(stepBody))
	if err != nil {
		t.Fatalf("step session: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, resp.StatusCode)
	}
	var stepResp map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&stepResp); err != nil {
		t.Fatalf("decode step response: %v", err)
	}
	if stepResp["stdout"] == "" {
		t.Fatalf("expected stdout")
	}
	if stepResp["stderr"] == "" {
		t.Fatalf("expected stderr")
	}
	if stepResp["status"] != "accepted" {
		t.Fatalf("expected accepted status")
	}
	_ = resp.Body.Close()

	resp, err = http.Post(server.URL+"/sessions/session-1/terminate", "application/json", bytes.NewReader([]byte("{}")))
	if err != nil {
		t.Fatalf("terminate session: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, resp.StatusCode)
	}
	_ = resp.Body.Close()
}
