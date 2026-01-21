package unit

import (
	"testing"

	"session-agent/internal/runtime"
)

func TestRunnerEnsureSessionReusesState(t *testing.T) {
	runner := runtime.NewRunner()
	first := runner.EnsureSession("session-1", "python")
	second := runner.EnsureSession("session-1", "node")

	if first != second {
		t.Fatalf("expected same session instance")
	}
	if first.Runtime != "python" {
		t.Fatalf("expected runtime to remain python, got %q", first.Runtime)
	}
}

func TestRunnerRemoveSession(t *testing.T) {
	runner := runtime.NewRunner()
	runner.EnsureSession("session-2", "python")
	if _, ok := runner.GetSession("session-2"); !ok {
		t.Fatalf("expected session to exist")
	}
	runner.RemoveSession("session-2")
	if _, ok := runner.GetSession("session-2"); ok {
		t.Fatalf("expected session to be removed")
	}
}
