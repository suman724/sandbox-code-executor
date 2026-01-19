package sessions

import "context"

type TTLWorker struct {}

func (TTLWorker) Run(ctx context.Context) error {
	_ = ctx
	return nil
}
