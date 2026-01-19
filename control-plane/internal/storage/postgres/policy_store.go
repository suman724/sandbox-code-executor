package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"control-plane/internal/storage"
)

type PolicyStore struct {
	Pool *pgxpool.Pool
}

type AuditStore struct {
	Pool *pgxpool.Pool
}

func (s PolicyStore) Upsert(ctx context.Context, policy storage.Policy) error {
	if s.Pool == nil {
		return errors.New("nil pool")
	}
	_, err := s.Pool.Exec(ctx, `insert into policies (id, version) values ($1, $2) on conflict (id) do update set version = excluded.version`, policy.ID, policy.Version)
	return err
}

func (s AuditStore) Append(ctx context.Context, event storage.AuditEvent) error {
	if s.Pool == nil {
		return errors.New("nil pool")
	}
	_, err := s.Pool.Exec(ctx, `insert into audit_events (id, action, outcome) values ($1, $2, $3)`, event.ID, event.Action, event.Outcome)
	return err
}

func (s AuditStore) List(ctx context.Context) ([]storage.AuditEvent, error) {
	if s.Pool == nil {
		return nil, errors.New("nil pool")
	}
	rows, err := s.Pool.Query(ctx, `select id, action, outcome from audit_events order by id desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []storage.AuditEvent
	for rows.Next() {
		var event storage.AuditEvent
		if err := rows.Scan(&event.ID, &event.Action, &event.Outcome); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
