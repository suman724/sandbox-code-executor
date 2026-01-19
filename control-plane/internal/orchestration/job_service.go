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
	if ok, err := s.Enforcer.Evaluate(ctx, job); err != nil {
		return "", err
	} else if !ok {
		return "", errors.New("policy denied job")
	}
	if err := s.Store.Create(ctx, storage.Job{ID: job.ID, Status: string(job.Status)}); err != nil {
		return "", err
	}
	resp, err := s.Client.StartRun(ctx, client.RunRequest{
		JobID:     job.ID,
		Language:  job.Language,
		Code:      job.Code,
		Workspace: job.Workspace,
	})
	if err != nil {
		return "", err
	}
	return resp.RunID, nil
}
