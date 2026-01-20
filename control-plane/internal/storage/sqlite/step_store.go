package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"control-plane/internal/storage"
)

type SessionStepStore struct {
	DB *sql.DB
}

func (s SessionStepStore) Append(ctx context.Context, step storage.SessionStep) error {
	if s.DB == nil {
		return errors.New("nil db")
	}
	_, err := s.DB.ExecContext(ctx, `insert into session_steps (id, session_id, command, status) values (?, ?, ?, ?)`, step.ID, step.SessionID, step.Command, step.Status)
	return err
}

func (s SessionStepStore) List(ctx context.Context, sessionID string) ([]storage.SessionStep, error) {
	if s.DB == nil {
		return nil, errors.New("nil db")
	}
	rows, err := s.DB.QueryContext(ctx, `select id, session_id, command, status from session_steps where session_id = ? order by rowid`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var steps []storage.SessionStep
	for rows.Next() {
		var step storage.SessionStep
		if err := rows.Scan(&step.ID, &step.SessionID, &step.Command, &step.Status); err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return steps, nil
}
