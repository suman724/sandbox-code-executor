package runtime

import (
	"context"
	"errors"
)

type ServiceRunner struct{}

func (ServiceRunner) Start(ctx context.Context, serviceID string) error {
	_ = ctx
	if serviceID == "" {
		return errors.New("missing service id")
	}
	return nil
}

func (ServiceRunner) Stop(ctx context.Context, serviceID string) error {
	_ = ctx
	if serviceID == "" {
		return errors.New("missing service id")
	}
	return nil
}
