package execution

import (
	"context"
	"testing"

	"data-plane/internal/runtime"
)

type mockAdapter struct {
	seen string
}

func (m *mockAdapter) Run(code string) error {
	m.seen = code
	return nil
}

func TestRunnerExecutesCode(t *testing.T) {
	adapter := &mockAdapter{}
	reg := runtime.NewRegistry()
	reg.Register("go", adapter)
	runner := Runner{Registry: reg}
	_, err := runner.Run(context.Background(), "job-1", "go", "print")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if adapter.seen == "" {
		t.Fatalf("expected adapter to receive code")
	}
}
