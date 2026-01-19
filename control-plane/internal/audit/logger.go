package audit

import (
	"context"
	"time"
)

type Event struct {
	TenantID string
	ActorID  string
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

func (StdoutLogger) ServiceStarted(ctx context.Context, tenantID string, serviceID string) error {
	return StdoutLogger{}.Log(ctx, Event{
		TenantID: tenantID,
		Action:   "service_started",
		Outcome:  "ok",
		Time:     time.Now(),
		Detail:   serviceID,
	})
}
