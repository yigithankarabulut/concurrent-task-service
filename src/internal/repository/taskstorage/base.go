package taskstorage

import (
	"database/sql"
	. "github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
)

type TaskStorer interface {
	Set(Task) error
	Get(uint) (Task, error)
	Update(Task) error
	Delete(uint) error
	List(string) ([]Task, error)
}

type taskStorage struct {
	db *sql.DB
}

type TaskStorageOption func(*taskStorage)

func WithTaskDB(db *sql.DB) TaskStorageOption {
	return func(s *taskStorage) {
		s.db = db
	}
}

func NewTaskStorage(opts ...TaskStorageOption) TaskStorer {
	s := &taskStorage{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
