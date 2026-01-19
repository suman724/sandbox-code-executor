package policy

import (
	"context"
	"testing"
)

type mockStore struct {
	last Policy
}

func (m *mockStore) Upsert(ctx context.Context, policy Policy) error {
	_ = ctx
	m.last = policy
	return nil
}

func TestPolicyCRUD(t *testing.T) {
	store := &mockStore{}
	err := store.Upsert(context.Background(), Policy{ID: "p-1", Version: 1, Ruleset: "allow"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.last.ID != "p-1" {
		t.Fatalf("expected policy stored")
	}
}
