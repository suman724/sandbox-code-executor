package orchestration

import "context"

type DegradationMode string

const (
	DegradationNone     DegradationMode = "none"
	DegradationReadOnly DegradationMode = "read_only"
)

type DegradationController interface {
	Mode(ctx context.Context) DegradationMode
}

type StaticDegradationController struct {
	mode DegradationMode
}

func (s StaticDegradationController) Mode(ctx context.Context) DegradationMode {
	_ = ctx
	return s.mode
}
