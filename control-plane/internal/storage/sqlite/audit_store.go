package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"control-plane/internal/storage"
)

type AuditStore struct {
	DB *sql.DB
}

func (s AuditStore) Append(ctx context.Context, event storage.AuditEvent) error {
	if s.DB == nil {
		return errors.New("nil db")
	}
	_, err := s.DB.ExecContext(ctx, `insert into audit_events (id, action, outcome) values (?, ?, ?)`, event.ID, event.Action, event.Outcome)
	return err
}

func (s AuditStore) List(ctx context.Context) ([]storage.AuditEvent, error) {
	if s.DB == nil {
		return nil, errors.New("nil db")
	}
	rows, err := s.DB.QueryContext(ctx, `select id, action, outcome from audit_events order by id desc`)
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
