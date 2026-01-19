package orchestration

import "context"

type WorkflowService struct{}

func (WorkflowService) Start(ctx context.Context, wf Workflow) error {
	_ = ctx
	_ = wf
	return nil
}
