package execution

import (
	"context"
	"testing"

	"data-plane/internal/runtime"
)

type mockAdapter struct {
	seen string
	err  error
}

func (m *mockAdapter) Run(code string) error {
	m.seen = code
	return m.err
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

func TestRunnerUnsupportedLanguage(t *testing.T) {
	runner := Runner{Registry: runtime.NewRegistry()}
	_, err := runner.Run(context.Background(), "job-1", "missing", "print")
	if err == nil {
		t.Fatalf("expected unsupported language error")
	}
}

func TestRunnerAdapterFailure(t *testing.T) {
	adapter := &mockAdapter{err: errTest}
	reg := runtime.NewRegistry()
	reg.Register("go", adapter)
	runner := Runner{Registry: reg}
	_, err := runner.Run(context.Background(), "job-1", "go", "print")
	if err == nil {
		t.Fatalf("expected adapter error")
	}
}

var errTest = runtimeError("adapter failed")

type runtimeError string

func (e runtimeError) Error() string { return string(e) }
