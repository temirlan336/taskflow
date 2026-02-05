package service

import (
	"context"
	"taskflow/internal/domain"
	"taskflow/internal/repository"
)

type TaskService struct {
	taskRepo  repository.TaskRepository
	eventRepo repository.EventRepository
	txManager repository.TxManager
}

func NewTaskService(r repository.TaskRepository, e repository.EventRepository, tx repository.TxManager) *TaskService {
	return &TaskService{taskRepo: r, eventRepo: e, txManager: tx}
}
func (t *TaskService) CreateTask(ctx context.Context, title string) (domain.Task, error) {
	tx, taskRepotx, eventRepoTx, err := t.txManager.Begin(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	task, err := taskRepotx.Create(ctx, title)
	if err != nil {
		return task, err
	}
	if err := eventRepoTx.Create(ctx, task.ID, "task_created"); err != nil {
		return task, err
	}
	if err := tx.Commit(); err != nil {
		return task, err
	}
	return task, nil
}

func (t *TaskService) GetTasks(ctx context.Context) ([]domain.Task, error) {
	return t.taskRepo.GetAll(ctx)
}
func (t *TaskService) GetTaskByID(ctx context.Context, id int) (domain.Task, error) {
	return t.taskRepo.GetByID(ctx, id)
}
func (t *TaskService) UpdateTask(ctx context.Context, id int, title string, completed bool) (domain.Task, error) {
	tx, taskRepotx, eventRepoTx, err := t.txManager.Begin(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer tx.Rollback()

	task, err := taskRepotx.Update(ctx, id, title, completed)
	if err != nil {
		return task, err
	}
	if err := eventRepoTx.Create(ctx, task.ID, "task_updated"); err != nil {
		return task, err
	}
	if err := tx.Commit(); err != nil {
		return task, err
	}
	return task, nil

}
func (t *TaskService) DeleteTask(ctx context.Context, id int) error {
	tx, taskRepotx, eventRepoTx, err := t.txManager.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := eventRepoTx.Create(ctx, id, "task_deleted"); err != nil {
		return err
	}

	if err := taskRepotx.Delete(ctx, id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
