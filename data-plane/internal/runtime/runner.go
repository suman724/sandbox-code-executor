package runtime

import "context"

type Runner interface {
	Run(ctx context.Context, jobID string, language string, code string) (string, error)
}
