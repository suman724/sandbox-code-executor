package sessions

import (
	"context"
	"testing"
	"time"

	"control-plane/internal/orchestration"
	"control-plane/internal/policy"
	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

type mockSessionStore struct{}

func (mockSessionStore) Create(ctx context.Context, session storage.Session) error {
	_ = ctx
	if session.ID == "" {
		return storageError("missing id")
	}
	return nil
}

type storageError string

func (e storageError) Error() string { return string(e) }

type mockEvaluator struct {
	allowed bool
}

func (m mockEvaluator) Evaluate(ctx context.Context, input any) (policy.Decision, error) {
	_ = ctx
	_ = input
	return policy.Decision{Allowed: m.allowed}, nil
}

func TestSessionLifecycle(t *testing.T) {
	svc := Service{
		Store:  mockSessionStore{},
		Client: client.DataPlaneClient{},
		Enforcer: orchestration.PolicyEnforcer{
			Evaluator: mockEvaluator{allowed: true},
		},
	}
	_, err := svc.CreateSession(context.Background(), Session{ID: "s-1", TTL: time.Minute, Status: StatusActive})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSessionPolicyDenied(t *testing.T) {
	svc := Service{
		Store:  mockSessionStore{},
		Client: client.DataPlaneClient{},
		Enforcer: orchestration.PolicyEnforcer{
			Evaluator: mockEvaluator{allowed: false},
		},
	}
	_, err := svc.CreateSession(context.Background(), Session{ID: "s-1", TTL: time.Minute, Status: StatusActive})
	if err == nil {
		t.Fatalf("expected policy denial error")
	}
}
