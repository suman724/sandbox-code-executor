package audit

import (
	"context"
	"errors"
	"log"
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
	log.Printf("audit event action=%s outcome=%s tenant=%s actor=%s detail=%s", event.Action, event.Outcome, event.TenantID, event.ActorID, event.Detail)
	return nil
}

type StoreLogger struct {
	Store Store
}

func (l StoreLogger) Log(ctx context.Context, event Event) error {
	if l.Store == nil {
		return errors.New("missing audit store")
	}
	if err := l.Store.Append(ctx, event); err != nil {
		return err
	}
	return StdoutLogger{}.Log(ctx, event)
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

func (StdoutLogger) ServiceStopped(ctx context.Context, tenantID string, serviceID string) error {
	return StdoutLogger{}.Log(ctx, Event{
		TenantID: tenantID,
		Action:   "service_stopped",
		Outcome:  "ok",
		Time:     time.Now(),
		Detail:   serviceID,
	})
}

func (StdoutLogger) JobAccepted(ctx context.Context, tenantID string, jobID string) error {
	return StdoutLogger{}.Log(ctx, Event{
		TenantID: tenantID,
		Action:   "job_accepted",
		Outcome:  "ok",
		Time:     time.Now(),
		Detail:   jobID,
	})
}

func (StdoutLogger) JobRunning(ctx context.Context, tenantID string, jobID string) error {
	return StdoutLogger{}.Log(ctx, Event{
		TenantID: tenantID,
		Action:   "job_running",
		Outcome:  "ok",
		Time:     time.Now(),
		Detail:   jobID,
	})
}

func (StdoutLogger) JobFinished(ctx context.Context, tenantID string, jobID string, outcome string) error {
	return StdoutLogger{}.Log(ctx, Event{
		TenantID: tenantID,
		Action:   "job_finished",
		Outcome:  outcome,
		Time:     time.Now(),
		Detail:   jobID,
	})
}

func (StdoutLogger) SessionCreated(ctx context.Context, tenantID string, sessionID string) error {
	return StdoutLogger{}.Log(ctx, Event{
		TenantID: tenantID,
		Action:   "session_created",
		Outcome:  "ok",
		Time:     time.Now(),
		Detail:   sessionID,
	})
}

func (StdoutLogger) SessionStepAccepted(ctx context.Context, tenantID string, stepID string) error {
	return StdoutLogger{}.Log(ctx, Event{
		TenantID: tenantID,
		Action:   "session_step_accepted",
		Outcome:  "ok",
		Time:     time.Now(),
		Detail:   stepID,
	})
}
