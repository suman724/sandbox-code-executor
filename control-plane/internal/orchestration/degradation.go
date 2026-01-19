package orchestration

import (
	"context"
	"errors"
)

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

var ErrReadOnlyMode = errors.New("service in read-only mode")

func RequireWriteAllowed(ctx context.Context, controller DegradationController) error {
	if controller == nil {
		return nil
	}
	if controller.Mode(ctx) == DegradationReadOnly {
		return ErrReadOnlyMode
	}
	return nil
}
