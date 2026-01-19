package orchestration

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"control-plane/internal/audit"
)

type WorkflowStore interface {
	Create(ctx context.Context, workflow Workflow) error
	UpdateStatus(ctx context.Context, id string, status WorkflowStatus) error
}

type SharedMemoryStore interface {
	Put(ctx context.Context, workflowID string, key string, value string) error
	Get(ctx context.Context, workflowID string, key string) (string, bool)
}

type WorkflowStepRunner interface {
	RunStep(ctx context.Context, workflowID string, step WorkflowStep, memory SharedMemoryStore) (string, error)
}

type WorkflowService struct {
	Store  WorkflowStore
	Runner WorkflowStepRunner
	Memory SharedMemoryStore
	Logger audit.Logger
	Now    func() time.Time
}

func (s WorkflowService) Start(ctx context.Context, wf Workflow) error {
	if wf.ID == "" {
		return errors.New("missing workflow id")
	}
	if wf.TenantID == "" {
		return errors.New("missing tenant id")
	}
	if len(wf.Steps) == 0 {
		return errors.New("missing workflow steps")
	}
	if s.Store == nil {
		return errors.New("missing workflow store")
	}
	if s.Runner == nil {
		return errors.New("missing workflow runner")
	}
	now := time.Now
	if s.Now != nil {
		now = s.Now
	}
	memory := s.Memory
	if memory == nil {
		memory = NewMemoryStore()
	}
	steps, err := normalizeSteps(wf)
	if err != nil {
		return err
	}
	wf.Status = WorkflowRunning
	if wf.CreatedAt.IsZero() {
		wf.CreatedAt = now()
	}
	if wf.StartedAt.IsZero() {
		wf.StartedAt = now()
	}
	wf.Steps = steps
	if err := s.Store.Create(ctx, wf); err != nil {
		return err
	}
	if err := s.Store.UpdateStatus(ctx, wf.ID, WorkflowRunning); err != nil {
		return err
	}
	if s.Logger != nil {
		_ = audit.WorkflowStarted(ctx, s.Logger, wf.TenantID, wf.ID)
	}
	for i := range wf.Steps {
		step := wf.Steps[i]
		step.Status = WorkflowStepRunning
		step.StartedAt = now()
		if s.Logger != nil {
			_ = audit.WorkflowStepStarted(ctx, s.Logger, wf.TenantID, wf.ID, step.ID)
		}
		jobID, err := s.Runner.RunStep(ctx, wf.ID, step, memory)
		if err != nil {
			step.Status = WorkflowStepFailed
			step.EndedAt = now()
			_ = s.Store.UpdateStatus(ctx, wf.ID, WorkflowFailed)
			if s.Logger != nil {
				_ = audit.WorkflowStepFinished(ctx, s.Logger, wf.TenantID, wf.ID, step.ID, "failed")
			}
			return err
		}
		step.JobID = jobID
		step.Status = WorkflowStepSucceeded
		step.EndedAt = now()
		if s.Logger != nil {
			_ = audit.WorkflowStepFinished(ctx, s.Logger, wf.TenantID, wf.ID, step.ID, "succeeded")
		}
	}
	wf.Status = WorkflowFinished
	wf.CompletedAt = now()
	if err := s.Store.UpdateStatus(ctx, wf.ID, wf.Status); err != nil {
		return err
	}
	if s.Logger != nil {
		_ = audit.WorkflowFinished(ctx, s.Logger, wf.TenantID, wf.ID)
	}
	return nil
}

func normalizeSteps(wf Workflow) ([]WorkflowStep, error) {
	steps := make([]WorkflowStep, len(wf.Steps))
	copy(steps, wf.Steps)
	seen := map[int]struct{}{}
	for i := range steps {
		if steps[i].AgentID == "" {
			return nil, errors.New("missing workflow step agent id")
		}
		if steps[i].Sequence <= 0 {
			steps[i].Sequence = i + 1
		}
		if _, exists := seen[steps[i].Sequence]; exists {
			return nil, errors.New("duplicate workflow step sequence")
		}
		seen[steps[i].Sequence] = struct{}{}
		if steps[i].WorkflowID == "" {
			steps[i].WorkflowID = wf.ID
		}
		if steps[i].ID == "" {
			steps[i].ID = fmt.Sprintf("%s-step-%d", wf.ID, steps[i].Sequence)
		}
		if steps[i].Status == "" {
			steps[i].Status = WorkflowStepQueued
		}
	}
	sort.Slice(steps, func(i, j int) bool {
		return steps[i].Sequence < steps[j].Sequence
	})
	return steps, nil
}
