package orchestration

import (
	"context"
	"errors"

	"control-plane/internal/storage"
	"control-plane/pkg/client"
)

type JobService struct {
	Store    storage.JobStore
	Client   client.DataPlaneClient
	Enforcer PolicyEnforcer
}

func (s JobService) CreateJob(ctx context.Context, job Job) (string, error) {
	if job.ID == "" {
		return "", errors.New("missing job id")
	}
	if job.Status == "" {
		job.Status = JobQueued
	}
	if ok, err := s.Enforcer.Evaluate(ctx, job); err != nil {
		return "", err
	} else if !ok {
		return "", errors.New("policy denied job")
	}
	if err := s.Store.Create(ctx, storage.Job{ID: job.ID, Status: string(job.Status)}); err != nil {
		return "", err
	}
	resp, err := s.Client.StartRun(ctx, client.RunRequest{
		JobID:        job.ID,
		PolicyID:     job.PolicyID,
		Language:     job.Language,
		Code:         job.Code,
		WorkspaceRef: job.Workspace,
	})
	if err != nil {
		_ = s.Store.UpdateStatus(ctx, job.ID, string(JobFailed))
		return "", err
	}
	if err := s.Store.UpdateStatus(ctx, job.ID, string(JobRunning)); err != nil {
		return "", err
	}
	return resp.RunID, nil
}
