package integration

import (
	"testing"

	"data-plane/internal/runtime"
)

func TestSessionAgentLocalRuntimePersistence(t *testing.T) {
	runtimeRunner := runtime.NewLocalSessionRuntime()
	route, err := runtimeRunner.StartSession(nil, "session-local", "", "", "python")
	if err != nil {
		t.Fatalf("start session: %v", err)
	}
	_, err = runtimeRunner.RunStep(nil, route.RuntimeID, "x = 10")
	if err != nil {
		t.Fatalf("run step 1: %v", err)
	}
	output, err := runtimeRunner.RunStep(nil, route.RuntimeID, "print(x + 5)")
	if err != nil {
		t.Fatalf("run step 2: %v", err)
	}
	if output.Stdout == "" {
		t.Fatalf("expected stdout")
	}
	if err := runtimeRunner.TerminateSession(nil, route.RuntimeID); err != nil {
		t.Fatalf("terminate session: %v", err)
	}
}
