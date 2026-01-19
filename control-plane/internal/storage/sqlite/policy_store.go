package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"control-plane/internal/storage"
)

type PolicyStore struct {
	DB *sql.DB
}

func (s PolicyStore) Upsert(ctx context.Context, policy storage.Policy) error {
	if s.DB == nil {
		return errors.New("nil db")
	}
	_, err := s.DB.ExecContext(ctx, `insert into policies (id, version) values (?, ?) on conflict(id) do update set version = excluded.version`, policy.ID, policy.Version)
	return err
}
