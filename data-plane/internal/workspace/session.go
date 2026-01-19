package workspace

import (
	"errors"
	"os"
	"path/filepath"
)

type SessionWorkspace struct {
	SessionID string
	Path      string
}

func NewSessionWorkspace(root string, sessionID string) (SessionWorkspace, error) {
	if sessionID == "" {
		return SessionWorkspace{}, errors.New("missing session id")
	}
	if root == "" {
		return SessionWorkspace{}, errors.New("missing workspace root")
	}
	path := filepath.Join(root, sessionID)
	if err := os.MkdirAll(path, 0o750); err != nil {
		return SessionWorkspace{}, err
	}
	return SessionWorkspace{SessionID: sessionID, Path: path}, nil
}
