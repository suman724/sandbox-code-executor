package unit

import (
	"testing"

	"session-agent/internal/runtime"
	"shared/sessionagent"
)

func TestRunnerEnsureSessionReusesState(t *testing.T) {
	runner := runtime.NewRunner()
	first, err := runner.RegisterSession(sessionagent.SessionRegisterRequest{
		SessionID: "session-1",
		Runtime:   "python",
		Token:     "token-1",
	})
	if err != nil {
		t.Fatalf("register session: %v", err)
	}
	defer runner.RemoveSession("session-1")
	second, err := runner.RegisterSession(sessionagent.SessionRegisterRequest{
		SessionID: "session-1",
		Runtime:   "node",
		Token:     "token-2",
	})
	if err != nil {
		t.Fatalf("register session again: %v", err)
	}

	if first != second {
		t.Fatalf("expected same session instance")
	}
	if first.Runtime != "python" {
		t.Fatalf("expected runtime to remain python, got %q", first.Runtime)
	}
}

func TestRunnerRemoveSession(t *testing.T) {
	runner := runtime.NewRunner()
	if _, err := runner.RegisterSession(sessionagent.SessionRegisterRequest{
		SessionID: "session-2",
		Runtime:   "python",
		Token:     "token-2",
	}); err != nil {
		t.Fatalf("register session: %v", err)
	}
	if _, ok := runner.GetSession("session-2"); !ok {
		t.Fatalf("expected session to exist")
	}
	runner.RemoveSession("session-2")
	if _, ok := runner.GetSession("session-2"); ok {
		t.Fatalf("expected session to be removed")
	}
}
