package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"control-plane/internal/api/handlers"
	"control-plane/internal/mcp"
	"control-plane/internal/orchestration"
	"control-plane/internal/policy"
	"control-plane/internal/sessions"
	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

type mcpJobStore struct {
	job storage.Job
}

func (m *mcpJobStore) Create(ctx context.Context, job storage.Job) error {
	m.job = job
	return nil
}

func (m *mcpJobStore) Get(ctx context.Context, id string) (storage.Job, error) {
	_ = ctx
	return storage.Job{ID: id, Status: m.job.Status}, nil
}

func (m *mcpJobStore) UpdateStatus(ctx context.Context, id string, status string) error {
	_ = id
	m.job.Status = status
	return nil
}

type mcpSessionStore struct {
	session storage.Session
}

func (m *mcpSessionStore) Create(ctx context.Context, session storage.Session) error {
	m.session = session
	return nil
}

func (m *mcpSessionStore) Get(ctx context.Context, id string) (storage.Session, error) {
	_ = ctx
	return storage.Session{ID: id, Status: m.session.Status}, nil
}

func (m *mcpSessionStore) UpdateStatus(ctx context.Context, id string, status string) error {
	_ = id
	m.session.Status = status
	return nil
}

type mcpArtifactStore struct {
	artifacts map[string]storage.Artifact
}

func (m *mcpArtifactStore) Put(ctx context.Context, artifact storage.Artifact) error {
	_ = ctx
	if m.artifacts == nil {
		m.artifacts = map[string]storage.Artifact{}
	}
	m.artifacts[artifact.ID] = artifact
	return nil
}

func (m *mcpArtifactStore) Get(ctx context.Context, id string) (storage.Artifact, error) {
	_ = ctx
	if artifact, ok := m.artifacts[id]; ok {
		return artifact, nil
	}
	return storage.Artifact{}, nil
}

func (m *mcpArtifactStore) SignedDownloadURL(id string) (string, error) {
	return "https://example.test/artifacts/" + id, nil
}

type mcpAllowAllEvaluator struct{}

func (mcpAllowAllEvaluator) Evaluate(ctx context.Context, input any) (policy.Decision, error) {
	_ = ctx
	_ = input
	return policy.Decision{Allowed: true}, nil
}

type mcpStepRunner struct {
	stepID string
}

func (m mcpStepRunner) RunStep(ctx context.Context, sessionID string, command string) (sessions.StepResult, error) {
	_ = ctx
	_ = sessionID
	_ = command
	return sessions.StepResult{ID: m.stepID}, nil
}

type mcpWorkflowStore struct{}

func (mcpWorkflowStore) Create(ctx context.Context, workflow orchestration.Workflow) error {
	_ = ctx
	_ = workflow
	return nil
}

func (mcpWorkflowStore) UpdateStatus(ctx context.Context, id string, status orchestration.WorkflowStatus) error {
	_ = ctx
	_ = id
	_ = status
	return nil
}

type mcpWorkflowRunner struct{}

func (mcpWorkflowRunner) RunStep(ctx context.Context, workflowID string, step orchestration.WorkflowStep, memory orchestration.SharedMemoryStore) (string, error) {
	_ = ctx
	_ = workflowID
	_ = step
	_ = memory
	return "job-" + step.ID, nil
}

func TestMCPToolsIntegration(t *testing.T) {
	t.Setenv("AUTHZ_BYPASS", "true")
	dataPlane := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"run_id":"run-1"}`))
	}))
	t.Cleanup(dataPlane.Close)

	jobStore := &mcpJobStore{}
	jobService := orchestration.JobService{
		Store:  jobStore,
		Client: client.DataPlaneClient{BaseURL: dataPlane.URL, Client: dataPlane.Client()},
		Enforcer: orchestration.PolicyEnforcer{
			Evaluator: mcpAllowAllEvaluator{},
		},
	}
	sessionStore := &mcpSessionStore{}
	sessionService := sessions.Service{
		Store:  sessionStore,
		Client: client.DataPlaneClient{BaseURL: dataPlane.URL, Client: dataPlane.Client()},
		Enforcer: orchestration.PolicyEnforcer{
			Evaluator: mcpAllowAllEvaluator{},
		},
	}
	workflowService := orchestration.WorkflowService{
		Store:  mcpWorkflowStore{},
		Runner: mcpWorkflowRunner{},
	}

	router := mcp.RouterWithDependencies(mcp.Dependencies{
		JobsHandler: handlers.JobHandler{Service: jobService, Store: jobStore},
		SessionsHandler: handlers.SessionHandler{
			Service: sessionService,
			Stepper: sessions.StepService{Runner: mcpStepRunner{stepID: "step-1"}},
		},
		WorkflowsHandler: handlers.WorkflowHandler{Service: workflowService},
		ArtifactStore:    &mcpArtifactStore{},
	})

	jobPayload := map[string]any{
		"tenantId": "tenant-1",
		"agentId":  "agent-1",
		"policyId": "policy-1",
		"language": "python",
		"code":     "print('ok')",
	}
	jobBody, _ := json.Marshal(jobPayload)
	jobReq := httptest.NewRequest(http.MethodPost, "/tools/jobs", bytes.NewReader(jobBody))
	jobRec := httptest.NewRecorder()
	router.ServeHTTP(jobRec, jobReq)
	if jobRec.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, jobRec.Code)
	}

	sessionPayload := map[string]any{
		"tenantId":   "tenant-1",
		"agentId":    "agent-1",
		"policyId":   "policy-1",
		"runtime":    "python",
		"ttlSeconds": 60,
	}
	sessionBody, _ := json.Marshal(sessionPayload)
	sessionReq := httptest.NewRequest(http.MethodPost, "/tools/sessions", bytes.NewReader(sessionBody))
	sessionRec := httptest.NewRecorder()
	router.ServeHTTP(sessionRec, sessionReq)
	if sessionRec.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, sessionRec.Code)
	}

	stepReq := httptest.NewRequest(http.MethodPost, "/tools/sessions/session-1/steps", bytes.NewReader([]byte(`{"command":"ls"}`)))
	stepRec := httptest.NewRecorder()
	router.ServeHTTP(stepRec, stepReq)
	if stepRec.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, stepRec.Code)
	}

	workflowPayload := map[string]any{
		"tenantId": "tenant-1",
		"steps":    []string{"agent-1", "agent-2"},
	}
	workflowBody, _ := json.Marshal(workflowPayload)
	workflowReq := httptest.NewRequest(http.MethodPost, "/tools/workflows", bytes.NewReader(workflowBody))
	workflowRec := httptest.NewRecorder()
	router.ServeHTTP(workflowRec, workflowReq)
	if workflowRec.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, workflowRec.Code)
	}

	artifactPayload := map[string]any{
		"id":         "artifact-1",
		"name":       "stdout",
		"sizeBytes":  12,
		"checksum":   "sum",
		"storageUri": "s3://bucket/stdout",
	}
	artifactBody, _ := json.Marshal(artifactPayload)
	artifactReq := httptest.NewRequest(http.MethodPost, "/tools/artifacts/upload", bytes.NewReader(artifactBody))
	artifactRec := httptest.NewRecorder()
	router.ServeHTTP(artifactRec, artifactReq)
	if artifactRec.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, artifactRec.Code)
	}

	downloadReq := httptest.NewRequest(http.MethodGet, "/tools/artifacts/artifact-1/download", nil)
	downloadRec := httptest.NewRecorder()
	router.ServeHTTP(downloadRec, downloadReq)
	if downloadRec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, downloadRec.Code)
	}
}
