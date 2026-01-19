package execution

import (
	"context"
	"errors"

	"data-plane/internal/runtime"
)

type Runner struct {
	Registry runtime.Registry
	Deps     runtime.DependencyPolicy
}

func (r Runner) Run(ctx context.Context, jobID string, language string, code string) (string, error) {
	_ = ctx
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
