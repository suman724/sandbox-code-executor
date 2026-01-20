package sessions

import (
	"context"
	"errors"
	"time"

	"control-plane/internal/audit"
	"control-plane/internal/orchestration"
	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

type Service struct {
	Store    storage.SessionStore
	Client   client.DataPlaneClient
	Enforcer orchestration.PolicyEnforcer
	Logger   audit.Logger
}

type StepRunner interface {
	RunStep(ctx context.Context, sessionID string, command string) (string, error)
}

type StepStore interface {
	AppendStep(ctx context.Context, step SessionStep) error
	ListSteps(ctx context.Context, sessionID string) ([]SessionStep, error)
}

type StepService struct {
	Runner StepRunner
	Store  StepStore
	Logger audit.Logger
}

func (s Service) CreateSession(ctx context.Context, session Session) (string, error) {
	if session.ID == "" {
		return "", errors.New("missing session id")
	}
	if session.Status == "" {
		session.Status = StatusActive
	}
	if ok, err := s.Enforcer.Evaluate(ctx, session); err != nil {
		return "", err
	} else if !ok {
		return "", errors.New("policy denied session")
	}
	if session.ExpiresAt.IsZero() {
		session.ExpiresAt = sessionExpires(session, time.Now())
	}
	resp, err := s.Client.StartSession(ctx, client.SessionCreateRequest{
		SessionID:    session.ID,
		PolicyID:     session.PolicyID,
		WorkspaceRef: session.ID,
	})
	if err != nil {
		return "", err
	}
	session.RuntimeID = resp.RuntimeID
	if err := s.Store.Create(ctx, storage.Session{ID: session.ID, Status: string(session.Status), RuntimeID: session.RuntimeID}); err != nil {
		return "", err
	}
	if err := s.Store.UpdateStatus(ctx, session.ID, string(StatusActive)); err != nil {
		return "", err
	}
	if s.Logger != nil {
		_ = s.Logger.Log(ctx, audit.Event{
			TenantID: session.TenantID,
			Action:   "session_created",
			Outcome:  "ok",
			Time:     time.Now(),
			Detail:   session.ID,
		})
	}
	return resp.RuntimeID, nil
}

func sessionExpires(session Session, now time.Time) time.Time {
	ttl := session.TTL
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return now.Add(ttl)
}

func (s StepService) Run(ctx context.Context, sessionID string, command string) (string, error) {
	if sessionID == "" {
		return "", errors.New("missing session id")
	}
	if command == "" {
		return "", errors.New("missing command")
	}
	if s.Runner == nil {
		return "", errors.New("missing runner")
	}
	stepID, err := s.Runner.RunStep(ctx, sessionID, command)
	if err != nil {
		return "", err
	}
	if s.Store != nil {
		_ = s.Store.AppendStep(ctx, SessionStep{
			ID:        stepID,
			SessionID: sessionID,
			Command:   command,
			Status:    "accepted",
			StartedAt: time.Now(),
		})
	}
	if s.Logger != nil {
		_ = s.Logger.Log(ctx, audit.Event{
			Action:   "session_step_accepted",
			Outcome:  "ok",
			Time:     time.Now(),
			Detail:   stepID,
		})
	}
	return stepID, nil
}
