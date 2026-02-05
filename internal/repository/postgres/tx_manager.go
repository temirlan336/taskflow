package postgres

import (
	"context"
	"database/sql"
	"taskflow/internal/repository"
)

type TxManagerImpl struct {
	db *sql.DB
}

func NewTxManagerImpl(db *sql.DB) *TxManagerImpl {
	return &TxManagerImpl{db: db}
}

func (t *TxManagerImpl) Begin(ctx context.Context) (repository.Tx, repository.TaskRepository, repository.EventRepository, error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	taskRepoTx := NewTaskStorage(tx)
	eventRepoTx := NewEventStorage(tx)
	return tx, taskRepoTx, eventRepoTx, nil
}
