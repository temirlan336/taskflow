package memory

import (
	"context"
	"sync"
	"taskflow/internal/domain"
	"time"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	tasks  map[int]domain.Task
	nextID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks: make(map[int]domain.Task),
	}
}

func (s *MemoryStorage) Create(ctx context.Context, title string) (domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := domain.Task{
		ID:        s.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	s.tasks[s.nextID] = task
	s.nextID++

	return task, nil
}

func (s *MemoryStorage) GetAll(ctx context.Context) ([]domain.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result, nil
}

func (s *MemoryStorage) GetByID(ctx context.Context, id int) (domain.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return domain.Task{}, domain.ErrNotFound
	}
	return task, nil
}

func (s *MemoryStorage) Update(ctx context.Context, id int, title string, completed bool) (domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return domain.Task{}, domain.ErrNotFound
	}

	task.Title = title
	task.Completed = completed
	s.tasks[id] = task

	return task, nil
}

func (s *MemoryStorage) Delete(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return domain.ErrNotFound
	}

	delete(s.tasks, id)
	return nil
}
