package sessions

import (
	"context"
	"time"

	"control-plane/pkg/client"
)

type DataPlaneStepRunner struct {
	Client client.DataPlaneClient
}

func (r DataPlaneStepRunner) RunStep(ctx context.Context, sessionID string, command string) (string, error) {
	if err := r.Client.RunSessionStep(ctx, sessionID, command); err != nil {
		return "", err
	}
	return "step-" + time.Now().UTC().Format("20060102150405.000000000"), nil
}
