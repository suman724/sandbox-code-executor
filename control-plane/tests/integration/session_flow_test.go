package integration

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
	created []storage.Session
}

func (m *mockSessionStore) Create(ctx context.Context, session storage.Session) error {
	_ = ctx
	m.created = append(m.created, session)
	return nil
}

func (m *mockSessionStore) Get(ctx context.Context, id string) (storage.Session, error) {
	_ = ctx
	return storage.Session{ID: id, Status: string(sessions.StatusActive)}, nil
}

func (m *mockSessionStore) UpdateStatus(ctx context.Context, id string, status string) error {
	_ = id
	_ = status
	return nil
}

type allowAllSessionEvaluator struct{}

func (allowAllSessionEvaluator) Evaluate(ctx context.Context, input any) (policy.Decision, error) {
	_ = ctx
	_ = input
	return policy.Decision{Allowed: true}, nil
}

func TestSessionFlowIntegration(t *testing.T) {
	dataPlane := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"id":"session-1","runtimeId":"runtime-123","status":"running"}`))
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
	handler := handlers.SessionHandler{
		Service: service,
		Stepper: sessions.StepService{
			Runner: mockStepRunner{stepID: "step-1"},
			Store:  &mockStepStore{},
		},
	}

	payload := map[string]any{
		"tenantId":   "tenant-1",
		"agentId":    "agent-1",
		"policyId":   "policy-1",
		"runtime":    "python",
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
	if len(store.created) != 1 {
		t.Fatalf("expected session to be created")
	}

	stepReq := httptest.NewRequest(http.MethodPost, "/sessions/session-1/steps", bytes.NewReader([]byte(`{"command":"ls"}`)))
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("sessionId", "session-1")
	stepReq = stepReq.WithContext(context.WithValue(stepReq.Context(), chi.RouteCtxKey, routeCtx))
	stepRec := httptest.NewRecorder()
	handler.ServeHTTP(stepRec, stepReq)
	if stepRec.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, stepRec.Code)
	}
	var stepResp map[string]string
	if err := json.NewDecoder(stepRec.Body).Decode(&stepResp); err != nil {
		t.Fatalf("decode step response: %v", err)
	}
	if stepResp["stdout"] != "out" {
		t.Fatalf("expected stdout to be propagated")
	}
	stepStore := handler.Stepper.Store.(*mockStepStore)
	if len(stepStore.steps) != 1 {
		t.Fatalf("expected step to be stored")
	}
}

type mockStepRunner struct {
	stepID string
}

func (m mockStepRunner) RunStep(ctx context.Context, sessionID string, command string) (sessions.StepResult, error) {
	_ = ctx
	_ = sessionID
	_ = command
	return sessions.StepResult{ID: m.stepID, Stdout: "out", Stderr: ""}, nil
}

type mockStepStore struct {
	steps []sessions.SessionStep
}

func (m *mockStepStore) AppendStep(ctx context.Context, step sessions.SessionStep) error {
	_ = ctx
	m.steps = append(m.steps, step)
	return nil
}

func (m *mockStepStore) ListSteps(ctx context.Context, sessionID string) ([]sessions.SessionStep, error) {
	_ = ctx
	_ = sessionID
	return m.steps, nil
}
