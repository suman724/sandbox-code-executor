package runtime

import "context"

type SessionRuntime interface {
	StartSession(ctx context.Context, sessionID string, policyID string, workspaceRef string) (string, error)
	RunStep(ctx context.Context, runtimeID string, command string) error
	TerminateSession(ctx context.Context, runtimeID string) error
}
