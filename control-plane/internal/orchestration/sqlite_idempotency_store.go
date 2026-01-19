package orchestration

import (
	"context"
	"database/sql"
	"errors"
)

type SQLiteIdempotencyStore struct {
	DB *sql.DB
}

func (s *SQLiteIdempotencyStore) Get(ctx context.Context, key string) (string, bool, error) {
	if s.DB == nil {
		return "", false, errors.New("nil db")
	}
	var value string
	err := s.DB.QueryRowContext(ctx, `select value from idempotency_keys where key = ?`, key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}
	return value, true, nil
}

func (s *SQLiteIdempotencyStore) Put(ctx context.Context, key string, value string) error {
	if s.DB == nil {
		return errors.New("nil db")
	}
	_, err := s.DB.ExecContext(ctx, `insert into idempotency_keys (key, value) values (?, ?) on conflict(key) do nothing`, key, value)
	return err
}
