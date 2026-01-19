package execution

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"data-plane/internal/runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type Runner struct {
	Registry      runtime.Registry
	Deps          runtime.DependencyPolicy
	WorkspaceRoot string
}

var (
	runMetricsOnce      sync.Once
	runLatencyHistogram metric.Float64Histogram
	runDeniedCounter    metric.Int64Counter
)

func (r Runner) Run(ctx context.Context, jobID string, language string, code string) (string, error) {
	runMetricsOnce.Do(initRunMetrics)
	start := time.Now()
	defer func() {
		if runLatencyHistogram != nil {
			runLatencyHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))
		}
	}()

	if jobID == "" {
		return "", errors.New("missing job id")
	}
	if err := r.ensureWorkspace(jobID); err != nil {
		return "", err
	}
	if err := runtime.ValidateDependencies(r.Deps); err != nil {
		if runDeniedCounter != nil {
			runDeniedCounter.Add(ctx, 1)
		}
		return "", err
	}
	adapter, ok := r.Registry.Adapter(language)
	if !ok {
		if runDeniedCounter != nil {
			runDeniedCounter.Add(ctx, 1)
		}
		return "", errors.New("unsupported language")
	}
	if err := adapter.Run(code); err != nil {
		return "", err
	}
	return jobID + "-run", nil
}

func (r Runner) ensureWorkspace(jobID string) error {
	if r.WorkspaceRoot == "" {
		return nil
	}
	path := filepath.Join(r.WorkspaceRoot, jobID)
	return os.MkdirAll(path, 0o750)
}

func initRunMetrics() {
	meter := otel.Meter("data-plane.execution")
	runLatencyHistogram, _ = meter.Float64Histogram("dataplane.run.latency_ms")
	runDeniedCounter, _ = meter.Int64Counter("dataplane.run.denied")
}
