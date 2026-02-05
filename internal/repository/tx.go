package repository

import "context"

type Tx interface { // минимальный контракт на завершение транзакции
	Commit() error
	Rollback() error
}

type TxManager interface { // то, что умеет начать транзакцию и вернуть:
	Begin(ctx context.Context) (Tx, TaskRepository, EventRepository, error)
}
