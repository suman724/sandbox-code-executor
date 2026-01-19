package orchestration

import (
	"context"

	"control-plane/internal/policy"
)

type PolicyEnforcer struct {
	Evaluator policy.Evaluator
}

func (p PolicyEnforcer) Evaluate(ctx context.Context, input any) (bool, error) {
	decision, err := p.Evaluator.Evaluate(ctx, input)
	if err != nil {
		return false, err
	}
	return decision.Allowed, nil
}
