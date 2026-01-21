package integration

import (
	"context"
	"strings"
	"testing"

	"data-plane/internal/runtime"
)

func TestSessionIsolation(t *testing.T) {
	runtimeRunner := runtime.NewLocalSessionRuntime()
	routeA, err := runtimeRunner.StartSession(context.Background(), "session-a", "", "", "python")
	if err != nil {
		t.Fatalf("start session-a: %v", err)
	}
	routeB, err := runtimeRunner.StartSession(context.Background(), "session-b", "", "", "python")
	if err != nil {
		_ = runtimeRunner.TerminateSession(context.Background(), routeA.RuntimeID)
		t.Fatalf("start session-b: %v", err)
	}
	defer func() {
		if err := runtimeRunner.TerminateSession(context.Background(), routeA.RuntimeID); err != nil {
			t.Fatalf("terminate session-a: %v", err)
		}
		if err := runtimeRunner.TerminateSession(context.Background(), routeB.RuntimeID); err != nil {
			t.Fatalf("terminate session-b: %v", err)
		}
	}()

	if _, err := runtimeRunner.RunStep(context.Background(), routeA.RuntimeID, "x = 1"); err != nil {
		t.Fatalf("step session-a: %v", err)
	}
	if _, err := runtimeRunner.RunStep(context.Background(), routeB.RuntimeID, "x = 10"); err != nil {
		t.Fatalf("step session-b: %v", err)
	}

	outA, err := runtimeRunner.RunStep(context.Background(), routeA.RuntimeID, "print(x)")
	if err != nil {
		t.Fatalf("read session-a: %v", err)
	}
	outB, err := runtimeRunner.RunStep(context.Background(), routeB.RuntimeID, "print(x)")
	if err != nil {
		t.Fatalf("read session-b: %v", err)
	}

	if strings.TrimSpace(outA.Stdout) != "1" {
		t.Fatalf("expected session-a stdout 1, got %q", outA.Stdout)
	}
	if strings.TrimSpace(outB.Stdout) != "10" {
		t.Fatalf("expected session-b stdout 10, got %q", outB.Stdout)
	}
}
