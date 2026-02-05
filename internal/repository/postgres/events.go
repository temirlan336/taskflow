package postgres

import (
	"context"
)

type EventStorage struct {
	q Execer
	// db *sql.DB
}

func NewEventStorage(e Execer) *EventStorage {
	return &EventStorage{q: e}
}
func (event *EventStorage) Create(ctx context.Context, taskID int, eventType string) error {
	_, err := event.q.ExecContext(ctx, `INSERT INTO events (task_id, event_type) VALUES ($1, $2)`, taskID, eventType)

	if err != nil {
		return err
	}

	return nil
}
