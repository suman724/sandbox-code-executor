package execution

import "context"

type SessionRunner struct {
	Runner Runner
}

func (s SessionRunner) RunStep(ctx context.Context, sessionID string, command string) (string, error) {
	_ = ctx
	_ = sessionID
	_, err := s.Runner.Run(ctx, sessionID, "session", command)
	return "step-1", err
}
