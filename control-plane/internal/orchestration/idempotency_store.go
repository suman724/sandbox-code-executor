package orchestration

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IdempotencyStore interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Put(ctx context.Context, key string, value string) error
}

type InMemoryIdempotencyStore struct {
	items map[string]string
}

func NewInMemoryIdempotencyStore() *InMemoryIdempotencyStore {
	return &InMemoryIdempotencyStore{items: map[string]string{}}
}

func (s *InMemoryIdempotencyStore) Get(ctx context.Context, key string) (string, bool, error) {
	_ = ctx
	value, ok := s.items[key]
	return value, ok, nil
}

func (s *InMemoryIdempotencyStore) Put(ctx context.Context, key string, value string) error {
	_ = ctx
	s.items[key] = value
	return nil
}

type PostgresIdempotencyStore struct {
	Pool *pgxpool.Pool
}

func (s *PostgresIdempotencyStore) Get(ctx context.Context, key string) (string, bool, error) {
	if s.Pool == nil {
		return "", false, errors.New("nil pool")
	}
	var value string
	err := s.Pool.QueryRow(ctx, `select value from idempotency_keys where key = $1`, key).Scan(&value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}
	return value, true, nil
}

func (s *PostgresIdempotencyStore) Put(ctx context.Context, key string, value string) error {
	if s.Pool == nil {
		return errors.New("nil pool")
	}
	_, err := s.Pool.Exec(ctx, `insert into idempotency_keys (key, value) values ($1, $2) on conflict (key) do nothing`, key, value)
	return err
}

func ResolveIdempotency(ctx context.Context, store IdempotencyStore, key string, create func(context.Context) (string, error)) (string, error) {
	if store == nil || key == "" {
		return create(ctx)
	}
	if value, ok, err := store.Get(ctx, key); err != nil {
		return "", err
	} else if ok {
		return value, nil
	}
	value, err := create(ctx)
	if err != nil {
		return "", err
	}
	if err := store.Put(ctx, key, value); err != nil {
		return "", err
	}
	return value, nil
}
