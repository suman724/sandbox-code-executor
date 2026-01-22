package runtime

import "errors"

var (
	ErrRuntimeNotFound    = errors.New("runtime not found")
	ErrRuntimeUnavailable = errors.New("runtime unavailable")
)
