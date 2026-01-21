package integration

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"data-plane/internal/runtime"
)

func TestSessionConcurrency(t *testing.T) {
	runtimeRunner := runtime.NewLocalSessionRuntime()
	route, err := runtimeRunner.StartSession(context.Background(), "session-concurrency", "", "", "python")
	if err != nil {
		t.Fatalf("start session: %v", err)
	}
	defer func() {
		if err := runtimeRunner.TerminateSession(context.Background(), route.RuntimeID); err != nil {
			t.Fatalf("terminate session: %v", err)
		}
	}()

	const steps = 5
	var wg sync.WaitGroup
	errs := make(chan error, steps)

	for i := 0; i < steps; i++ {
		stepID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			code := fmt.Sprintf("print(%d)", stepID)
			output, err := runtimeRunner.RunStep(context.Background(), route.RuntimeID, code)
			if err != nil {
				errs <- err
				return
			}
			stdout := strings.TrimSpace(output.Stdout)
			if stdout != fmt.Sprintf("%d", stepID) {
				errs <- fmt.Errorf("unexpected stdout %q for step %d", stdout, stepID)
			}
		}()
	}

	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent step failed: %v", err)
		}
	}
}
