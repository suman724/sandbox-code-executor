package isolation

import "context"

type Limits struct {
	CPU    int
	Memory int
	Disk   int
	PIDs   int
}

func EnforceLimits(ctx context.Context, limits Limits) error {
	_ = ctx
	_ = limits
	return nil
}
