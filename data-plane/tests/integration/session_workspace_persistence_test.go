package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"data-plane/internal/runtime"
)

func TestSessionWorkspacePersistence(t *testing.T) {
	workspaceRoot := t.TempDir()
	runtimeRunner := runtime.NewLocalSessionRuntime()
	route, err := runtimeRunner.StartSession(context.Background(), "session-workspace", "", workspaceRoot, "python")
	if err != nil {
		t.Fatalf("start session: %v", err)
	}
	defer func() {
		if err := runtimeRunner.TerminateSession(context.Background(), route.RuntimeID); err != nil {
			t.Fatalf("terminate session: %v", err)
		}
	}()

	filePath := filepath.Join(workspaceRoot, "note.txt")
	writeCode := fmt.Sprintf("with open(%q, 'w') as f:\n    f.write('hello')", filePath)
	if _, err := runtimeRunner.RunStep(context.Background(), route.RuntimeID, writeCode); err != nil {
		t.Fatalf("write step: %v", err)
	}
	readCode := fmt.Sprintf("print(open(%q).read())", filePath)
	output, err := runtimeRunner.RunStep(context.Background(), route.RuntimeID, readCode)
	if err != nil {
		t.Fatalf("read step: %v", err)
	}
	if strings.TrimSpace(output.Stdout) != "hello" {
		t.Fatalf("expected workspace output hello, got %q", output.Stdout)
	}
}
