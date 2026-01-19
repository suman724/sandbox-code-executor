package isolation

import (
	"context"
	"errors"
)

type EgressPolicy struct {
	AllowList []string
	Requested []string
}

func EnforceEgress(ctx context.Context, policy EgressPolicy) error {
	_ = ctx
	if len(policy.Requested) == 0 {
		return nil
	}
	if len(policy.AllowList) == 0 {
		return errors.New("egress denied")
	}
	allowed := map[string]struct{}{}
	for _, host := range policy.AllowList {
		allowed[host] = struct{}{}
	}
	for _, host := range policy.Requested {
		if _, ok := allowed[host]; !ok {
			return errors.New("egress denied")
		}
	}
	return nil
}
