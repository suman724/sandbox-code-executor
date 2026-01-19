package runtime

import "context"

type ServiceRunner struct {}

func (ServiceRunner) Start(ctx context.Context, serviceID string) error {
	_ = ctx
	_ = serviceID
	return nil
}
