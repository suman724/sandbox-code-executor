package execution

import (
	"context"
	"time"
)

type SessionRunner struct {
	Runner Runner
}

func (s SessionRunner) RunStep(ctx context.Context, sessionID string, command string) (string, error) {
	_, err := s.Runner.Run(ctx, sessionID, "session", command)
	if err != nil {
		return "", err
	}
	return "step-" + time.Now().UTC().Format("20060102150405"), nil
}
