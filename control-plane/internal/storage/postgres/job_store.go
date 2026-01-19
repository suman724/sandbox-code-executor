package postgres

import (
	"context"

	"control-plane/internal/storage"
)

type JobStore struct{}

type SessionStore struct{}

func (JobStore) Create(ctx context.Context, job storage.Job) error {
	_ = ctx
	_ = job
	return nil
}

func (SessionStore) Create(ctx context.Context, session storage.Session) error {
	_ = ctx
	_ = session
	return nil
}
