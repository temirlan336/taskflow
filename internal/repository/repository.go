package repository

import (
	"context"
	"taskflow/internal/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, title string) (domain.Task, error)
	GetAll(ctx context.Context) ([]domain.Task, error)
	GetByID(ctx context.Context, id int) (domain.Task, error)
	Update(ctx context.Context, id int, title string, completed bool) (domain.Task, error)
	Delete(ctx context.Context, id int) error
}

type EventRepository interface {
	Create(ctx context.Context, taskID int, eventType string) error
}
