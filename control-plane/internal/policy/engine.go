package policy

import "context"

type Decision struct {
	Allowed bool
	Reason  string
}

type Evaluator interface {
	Evaluate(ctx context.Context, input any) (Decision, error)
}
