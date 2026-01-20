package factory

import (
	"context"
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"

	"github.com/jackc/pgx/v5/pgxpool"

	"control-plane/internal/storage"
	"control-plane/internal/storage/postgres"
	"control-plane/internal/storage/sqlite"
)

type StoreSet struct {
	JobStore         storage.JobStore
	SessionStore     storage.SessionStore
	SessionStepStore storage.SessionStepStore
	PolicyStore      storage.PolicyStore
	AuditStore       storage.AuditStore
	DB               *sql.DB
	Close            func() error
}

func NewStoreSet(ctx context.Context, driver string, dsn string) (StoreSet, error) {
	switch driver {
	case "postgres":
		if dsn == "" {
			return StoreSet{}, errors.New("missing postgres dsn")
		}
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			return StoreSet{}, err
		}
		return StoreSet{
			JobStore:         postgres.JobStore{Pool: pool},
			SessionStore:     postgres.SessionStore{Pool: pool},
			SessionStepStore: postgres.SessionStepStore{Pool: pool},
			PolicyStore:      postgres.PolicyStore{Pool: pool},
			AuditStore:       postgres.AuditStore{Pool: pool},
			Close: func() error {
				pool.Close()
				return nil
			},
		}, nil
	case "sqlite":
		if dsn == "" {
			return StoreSet{}, errors.New("missing sqlite dsn")
		}
		db, err := sql.Open("sqlite", dsn)
		if err != nil {
			return StoreSet{}, err
		}
		if err := sqlite.Bootstrap(db); err != nil {
			_ = db.Close()
			return StoreSet{}, err
		}
		return StoreSet{
			JobStore:         sqlite.JobStore{DB: db},
			SessionStore:     sqlite.SessionStore{DB: db},
			SessionStepStore: sqlite.SessionStepStore{DB: db},
			PolicyStore:      sqlite.PolicyStore{DB: db},
			AuditStore:       sqlite.AuditStore{DB: db},
			DB:               db,
			Close:            db.Close,
		}, nil
	default:
		return StoreSet{}, errors.New("unsupported database driver")
	}
}
