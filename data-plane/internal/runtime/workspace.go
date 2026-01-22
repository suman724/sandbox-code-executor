package runtime

import (
	"os"
	"path/filepath"

	"data-plane/internal/workspace"
)

func resolveWorkspaceDir(workspaceRef string, sessionID string) (string, error) {
	root := getenv("WORKSPACE_ROOT", "/tmp/sessions")
	ref := workspaceRef
	if ref == "" {
		ref = sessionID
	}
	if root == "" {
		return "", nil
	}
	ws, err := workspace.NewSessionWorkspace(root, ref)
	if err != nil {
		return "", err
	}
	return filepath.Clean(ws.Path), nil
}

func getenv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
