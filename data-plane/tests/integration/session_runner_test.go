package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
		"runtime":      "python",
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

	stepBody, err := json.Marshal(map[string]string{"command": "import sys\nprint('out')\nprint('err', file=sys.stderr)", "runtime": "python"})
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

	stepBody, err = json.Marshal(map[string]string{"command": "x = 41\nprint('set')", "runtime": "python"})
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
	resp.Body.Close()

	stepBody, err = json.Marshal(map[string]string{"command": "print(x + 1)", "runtime": "python"})
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
	stepResp = map[string]string{}
	if err := json.NewDecoder(resp.Body).Decode(&stepResp); err != nil {
		t.Fatalf("decode step response: %v", err)
	}
	if !strings.Contains(stepResp["stdout"], "42") {
		t.Fatalf("expected stdout to contain 42, got %q", stepResp["stdout"])
	}
	resp.Body.Close()

	resp, err = http.Post(server.URL+"/sessions/session-1/terminate", "application/json", bytes.NewReader([]byte("{}")))
	if err != nil {
		t.Fatalf("terminate session: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, resp.StatusCode)
	}
	_ = resp.Body.Close()
}
