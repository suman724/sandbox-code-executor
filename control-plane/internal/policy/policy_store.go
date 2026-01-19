package policy

import "context"

type Policy struct {
	ID      string
	Version int
	Ruleset string
}

type Store interface {
	Upsert(ctx context.Context, policy Policy) error
}
