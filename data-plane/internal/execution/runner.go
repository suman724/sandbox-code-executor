package execution

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"data-plane/internal/runtime"
)

type Runner struct {
	Registry      runtime.Registry
	Deps          runtime.DependencyPolicy
	WorkspaceRoot string
}

func (r Runner) Run(ctx context.Context, jobID string, language string, code string) (string, error) {
	if jobID == "" {
		return "", errors.New("missing job id")
	}
	if err := r.ensureWorkspace(jobID); err != nil {
		return "", err
	}
	if err := runtime.ValidateDependencies(r.Deps); err != nil {
		return "", err
	}
	adapter, ok := r.Registry.Adapter(language)
	if !ok {
		return "", errors.New("unsupported language")
	}
	if err := adapter.Run(code); err != nil {
		return "", err
	}
	return jobID + "-run", nil
}

func (r Runner) ensureWorkspace(jobID string) error {
	if r.WorkspaceRoot == "" {
		return nil
	}
	path := filepath.Join(r.WorkspaceRoot, jobID)
	return os.MkdirAll(path, 0o750)
}
