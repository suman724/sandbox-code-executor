package sessions

import (
	"context"
	"time"

	"control-plane/pkg/client"
)

type DataPlaneStepRunner struct {
	Client client.DataPlaneClient
}

func (r DataPlaneStepRunner) RunStep(ctx context.Context, sessionID string, command string) (StepResult, error) {
	resp, err := r.Client.RunSessionStep(ctx, sessionID, command)
	if err != nil {
		return StepResult{}, err
	}
	return StepResult{
		ID:     "step-" + time.Now().UTC().Format("20060102150405.000000000"),
		Stdout: resp.Stdout,
		Stderr: resp.Stderr,
	}, nil
}
