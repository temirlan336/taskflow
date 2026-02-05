package domain

import (
	"time"
)

type Event struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	EventType string    `json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
}
