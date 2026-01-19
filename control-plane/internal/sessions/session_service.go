package sessions

import (
	"context"
	"errors"
	"time"

	"control-plane/internal/orchestration"
	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

type Service struct {
	Store    storage.SessionStore
	Client   client.DataPlaneClient
	Enforcer orchestration.PolicyEnforcer
}

func (s Service) CreateSession(ctx context.Context, session Session) (string, error) {
	if ok, err := s.Enforcer.Evaluate(ctx, session); err != nil {
		return "", err
	} else if !ok {
		return "", errors.New("policy denied session")
	}
	if session.ExpiresAt.IsZero() {
		session.ExpiresAt = time.Now().Add(session.TTL)
	}
	if err := s.Store.Create(ctx, storage.Session{ID: session.ID, Status: string(session.Status)}); err != nil {
		return "", err
	}
	resp, err := s.Client.StartRun(ctx, client.RunRequest{
		JobID:        session.ID,
		PolicyID:     session.PolicyID,
		Language:     "session",
		Code:         "",
		WorkspaceRef: session.ID,
	})
	if err != nil {
		return "", err
	}
	return resp.RunID, nil
}
