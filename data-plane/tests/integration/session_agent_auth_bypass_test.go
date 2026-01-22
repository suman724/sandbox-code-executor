package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"data-plane/internal/runtime"
	"shared/sessionagent"
)

func TestAgentAuthBypassAcceptsSteps(t *testing.T) {
	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/steps" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"stepId":"step-1","status":"completed","stdout":"ok","stderr":""}`))
	}))
	defer agentServer.Close()

	client := runtime.NewAgentClient()
	result, err := client.RunStep(context.Background(), runtime.AgentRoute{
		Endpoint: agentServer.URL,
		AuthMode: "bypass",
	}, sessionagent.StepRequest{
		SessionID: "session-1",
		StepID:    "step-1",
		Code:      "print(1)",
		Runtime:   "python",
	})
	if err != nil {
		t.Fatalf("run step: %v", err)
	}
	if result.Stdout != "ok" {
		t.Fatalf("expected stdout ok, got %q", result.Stdout)
	}
}
