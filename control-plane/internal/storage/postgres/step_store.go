package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"control-plane/internal/storage"
)

type SessionStepStore struct {
	Pool *pgxpool.Pool
}

func (s SessionStepStore) Append(ctx context.Context, step storage.SessionStep) error {
	if s.Pool == nil {
		return errors.New("nil pool")
	}
	_, err := s.Pool.Exec(ctx, `insert into session_steps (id, session_id, command, status) values ($1, $2, $3, $4)`, step.ID, step.SessionID, step.Command, step.Status)
	return err
}

func (s SessionStepStore) List(ctx context.Context, sessionID string) ([]storage.SessionStep, error) {
	if s.Pool == nil {
		return nil, errors.New("nil pool")
	}
	rows, err := s.Pool.Query(ctx, `select id, session_id, command, status from session_steps where session_id = $1 order by id`, sessionID)
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
	if len(steps) == 0 {
		return nil, pgx.ErrNoRows
	}
	return steps, nil
}
