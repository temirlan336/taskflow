package postgres

import (
	"context"
	"database/sql"
	"taskflow/internal/domain"
)

type TaskStorage struct {
	q Execer
	// db *sql.DB
}

func NewTaskStorage(e Execer) *TaskStorage {
	return &TaskStorage{
		q: e,
	}
}

func (s *TaskStorage) Create(ctx context.Context, title string) (domain.Task, error) {
	var task domain.Task
	row := s.q.QueryRowContext(ctx, `INSERT INTO tasks (title) VALUES ($1) RETURNING id, title, completed, created_at`, title)

	err := row.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt)

	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (s *TaskStorage) GetAll(ctx context.Context) ([]domain.Task, error) {
	tasks := make([]domain.Task, 0)
	rows, err := s.q.QueryContext(ctx, `SELECT id, title, completed, created_at FROM tasks	ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil

}

func (s *TaskStorage) GetByID(ctx context.Context, id int) (domain.Task, error) {
	var task domain.Task
	row := s.q.QueryRowContext(ctx, `SELECT id, title, completed, created_at FROM tasks WHERE id = $1`, id)
	err := row.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt)

	if err == sql.ErrNoRows {
		return domain.Task{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *TaskStorage) Update(ctx context.Context, id int, title string, completed bool) (domain.Task, error) {
	var task domain.Task
	row := s.q.QueryRowContext(ctx, `UPDATE tasks SET title = $2, completed = $3 WHERE id = $1	RETURNING id, title, completed, created_at`, id, title, completed)
	err := row.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.Task{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (s *TaskStorage) Delete(ctx context.Context, id int) error {
	res, err := s.q.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
