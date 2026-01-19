package isolation

import "context"

type EgressPolicy struct {
	AllowList []string
}

func EnforceEgress(ctx context.Context, policy EgressPolicy) error {
	_ = ctx
	_ = policy
	return nil
}
