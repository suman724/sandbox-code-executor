package runtime

import "context"

type SessionRuntime interface {
	StartSession(ctx context.Context, sessionID string, policyID string, workspaceRef string, runtime string) (SessionRoute, error)
	RunStep(ctx context.Context, runtimeID string, command string) (StepOutput, error)
	TerminateSession(ctx context.Context, runtimeID string) error
}

type StepOutput struct {
	Stdout string
	Stderr string
}
