package orchestration

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"control-plane/internal/storage"
	"control-plane/pkg/client"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type JobService struct {
	Store    storage.JobStore
	Client   client.DataPlaneClient
	Enforcer PolicyEnforcer
}

var (
	jobMetricsOnce      sync.Once
	jobLatencyHistogram metric.Float64Histogram
	jobDeniedCounter    metric.Int64Counter
	jobQueueDepth       metric.Int64ObservableGauge
	jobQueued           atomic.Int64
)

func (s JobService) CreateJob(ctx context.Context, job Job) (string, error) {
	jobMetricsOnce.Do(initJobMetrics)
	start := time.Now()
	defer func() {
		if jobLatencyHistogram != nil {
			jobLatencyHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))
		}
	}()

	if job.ID == "" {
		return "", errors.New("missing job id")
	}
	if job.Status == "" {
		job.Status = JobQueued
	}
	if ok, err := s.Enforcer.Evaluate(ctx, job); err != nil {
		return "", err
	} else if !ok {
		if jobDeniedCounter != nil {
			jobDeniedCounter.Add(ctx, 1)
		}
		return "", errors.New("policy denied job")
	}
	if err := s.Store.Create(ctx, storage.Job{ID: job.ID, Status: string(job.Status)}); err != nil {
		return "", err
	}
	jobQueued.Add(1)
	resp, err := s.Client.StartRun(ctx, client.RunRequest{
		JobID:        job.ID,
		PolicyID:     job.PolicyID,
		Language:     job.Language,
		Code:         job.Code,
		WorkspaceRef: job.Workspace,
	})
	if err != nil {
		_ = s.Store.UpdateStatus(ctx, job.ID, string(JobFailed))
		jobQueued.Add(-1)
		return "", err
	}
	if err := s.Store.UpdateStatus(ctx, job.ID, string(JobRunning)); err != nil {
		jobQueued.Add(-1)
		return "", err
	}
	jobQueued.Add(-1)
	return resp.RunID, nil
}

func initJobMetrics() {
	meter := otel.Meter("control-plane.orchestration")
	jobLatencyHistogram, _ = meter.Float64Histogram("controlplane.jobs.latency_ms")
	jobDeniedCounter, _ = meter.Int64Counter("controlplane.jobs.policy_denied")
	jobQueueDepth, _ = meter.Int64ObservableGauge("controlplane.jobs.queue_depth")
	if jobQueueDepth != nil {
		_, _ = meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
			observer.ObserveInt64(jobQueueDepth, jobQueued.Load())
			return nil
		}, jobQueueDepth)
	}
}
