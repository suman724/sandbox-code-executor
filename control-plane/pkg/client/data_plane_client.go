package client

import "context"

type RunRequest struct {
	JobID     string
	Language  string
	Code      string
	Workspace string
}

type RunResponse struct {
	RunID string
}

type DataPlaneClient struct {
	BaseURL string
}

func (c DataPlaneClient) StartRun(ctx context.Context, req RunRequest) (RunResponse, error) {
	_ = ctx
	_ = req
	return RunResponse{}, nil
}
