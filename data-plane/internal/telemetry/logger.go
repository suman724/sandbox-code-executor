package telemetry

import (
	"context"
	"time"
)

type Event struct {
	TenantID string
	RunID    string
	Action   string
	Outcome  string
	Detail   string
	Time     time.Time
}

type Logger interface {
	Log(ctx context.Context, event Event) error
}

type StdoutLogger struct{}

func (StdoutLogger) Log(ctx context.Context, event Event) error {
	_ = ctx
	_ = event
	return nil
}

func (StdoutLogger) ServiceStopped(ctx context.Context, tenantID string, runID string) error {
	return StdoutLogger{}.Log(ctx, Event{
		TenantID: tenantID,
		RunID:    runID,
		Action:   "service_stopped",
		Outcome:  "ok",
		Time:     time.Now(),
		Detail:   runID,
	})
}
