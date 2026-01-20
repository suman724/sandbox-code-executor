package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"control-plane/internal/api/handlers"
	"control-plane/internal/orchestration"
	"control-plane/internal/policy"
	"control-plane/internal/sessions"
	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

type mockSessionStore struct {
	session storage.Session
}

func (m *mockSessionStore) Create(ctx context.Context, session storage.Session) error {
	m.session = session
	return nil
}

func (m *mockSessionStore) Get(ctx context.Context, id string) (storage.Session, error) {
	_ = ctx
	return m.session, nil
}

func (m *mockSessionStore) UpdateStatus(ctx context.Context, id string, status string) error {
	_ = id
	m.session.Status = status
	return nil
}

type allowAllSessionEvaluator struct{}

func (allowAllSessionEvaluator) Evaluate(ctx context.Context, input any) (policy.Decision, error) {
	_ = ctx
	_ = input
	return policy.Decision{Allowed: true}, nil
}

func TestSessionsContractCreate(t *testing.T) {
	dataPlane := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"id":"session-1","runtimeId":"runtime-1","status":"running"}`))
	}))
	t.Cleanup(dataPlane.Close)

	store := &mockSessionStore{}
	service := sessions.Service{
		Store:  store,
		Client: client.DataPlaneClient{BaseURL: dataPlane.URL, Client: dataPlane.Client()},
		Enforcer: orchestration.PolicyEnforcer{
			Evaluator: allowAllSessionEvaluator{},
		},
	}
	handler := handlers.SessionHandler{Service: service}

	payload := map[string]any{
		"tenantId":   "tenant-1",
		"agentId":    "agent-1",
		"policyId":   "policy-1",
		"ttlSeconds": 60,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/sessions", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestSessionsContractStep(t *testing.T) {
	handler := handlers.SessionHandler{
		Stepper: sessions.StepService{
			Runner: mockStepRunner{stepID: "step-1"},
		},
	}
	req := httptest.NewRequest(http.MethodPost, "/sessions/session-1/steps", bytes.NewReader([]byte(`{"command":"ls"}`)))
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("sessionId", "session-1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, rec.Code)
	}
	var stepResp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&stepResp); err != nil {
		t.Fatalf("decode step response: %v", err)
	}
	if stepResp["stdout"] != "ok" {
		t.Fatalf("expected stdout to be propagated")
	}
}

type mockStepRunner struct {
	stepID string
}

func (m mockStepRunner) RunStep(ctx context.Context, sessionID string, command string) (sessions.StepResult, error) {
	_ = ctx
	_ = sessionID
	_ = command
	return sessions.StepResult{ID: m.stepID, Stdout: "ok", Stderr: ""}, nil
}
