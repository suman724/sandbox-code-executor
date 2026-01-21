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

func TestSessionRoutingIsolation(t *testing.T) {
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

	sessionPayload := func(id string) []byte {
		payload, _ := json.Marshal(map[string]any{
			"sessionId":    id,
			"policyId":     "policy-1",
			"workspaceRef": id,
			"runtime":      "python",
		})
		return payload
	}

	resp, err := http.Post(server.URL+"/sessions", "application/json", bytes.NewReader(sessionPayload("session-1")))
	if err != nil {
		t.Fatalf("create session-1: %v", err)
	}
	resp.Body.Close()

	resp, err = http.Post(server.URL+"/sessions", "application/json", bytes.NewReader(sessionPayload("session-2")))
	if err != nil {
		t.Fatalf("create session-2: %v", err)
	}
	resp.Body.Close()

	step := func(id, code string) {
		body, _ := json.Marshal(map[string]string{"command": code, "runtime": "python"})
		resp, err := http.Post(server.URL+"/sessions/"+id+"/steps", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("step %s: %v", id, err)
		}
		resp.Body.Close()
	}

	step("session-1", "x = 1")
	step("session-2", "x = 10")

	body, _ := json.Marshal(map[string]string{"command": "print(x)", "runtime": "python"})
	resp, err = http.Post(server.URL+"/sessions/session-1/steps", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("read session-1: %v", err)
	}
	var output map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		t.Fatalf("decode session-1: %v", err)
	}
	resp.Body.Close()
	if output["stdout"] == "" {
		t.Fatalf("expected output for session-1")
	}

	resp, err = http.Post(server.URL+"/sessions/session-2/steps", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("read session-2: %v", err)
	}
	output = map[string]string{}
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		t.Fatalf("decode session-2: %v", err)
	}
	resp.Body.Close()
	if output["stdout"] == "" {
		t.Fatalf("expected output for session-2")
	}
}
