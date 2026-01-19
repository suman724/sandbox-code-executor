package contracts

import "fmt"

type ErrorCode string

const (
	ErrUnauthorized ErrorCode = "unauthorized"
	ErrForbidden    ErrorCode = "forbidden"
	ErrInvalidInput ErrorCode = "invalid_input"
	ErrNotFound     ErrorCode = "not_found"
	ErrInternal     ErrorCode = "internal"
)

type Error struct {
	Code    ErrorCode
	Message string
}

func (e Error) Error() string {
	if e.Message == "" {
		return string(e.Code)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
