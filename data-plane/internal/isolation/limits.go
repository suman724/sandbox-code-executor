package isolation

import (
	"context"
	"errors"
)

type Limits struct {
	CPU    int
	Memory int
	Disk   int
	PIDs   int
}

func EnforceLimits(ctx context.Context, limits Limits) error {
	_ = ctx
	if limits.CPU <= 0 {
		return errors.New("cpu limit required")
	}
	if limits.Memory <= 0 {
		return errors.New("memory limit required")
	}
	if limits.Disk <= 0 {
		return errors.New("disk limit required")
	}
	if limits.PIDs <= 0 {
		return errors.New("pid limit required")
	}
	return nil
}
