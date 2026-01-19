package postgres

import (
	"context"

	"control-plane/internal/storage"
)

type PolicyStore struct{}

type AuditStore struct{}

func (PolicyStore) Upsert(ctx context.Context, policy storage.Policy) error {
	_ = ctx
	_ = policy
	return nil
}

func (AuditStore) Append(ctx context.Context, event storage.AuditEvent) error {
	_ = ctx
	_ = event
	return nil
}
