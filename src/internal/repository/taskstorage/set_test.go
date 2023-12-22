package taskstorage_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
	. "github.com/yigithankarabulut/ConcurrentTaskService/src/internal/repository/taskstorage"
	"testing"
)

func Test_taskStorage_Set(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockStorage := NewTaskStorage(WithTaskDB(db))
	tasks := []models.Task{
		{
			ID:          1,
			Title:       "title",
			Description: "description",
			Status:      "status",
		},
		{
			ID:          2,
			Title:       "test",
			Description: "test",
			Status:      "test",
		},
	}
	for _, task := range tasks {
		mock.ExpectExec("INSERT INTO tasks").
			WithArgs(task.ID, task.Title, task.Description, task.Status).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	tests := []struct {
		name    string
		args    models.Task
		wantErr bool
	}{
		{
			name:    "Valid Task 1: No error is expected",
			args:    models.Task{ID: 1, Title: "title", Description: "description", Status: "status"},
			wantErr: false,
		},
		{
			name:    "Valid Task 2: No error is expected",
			args:    models.Task{ID: 2, Title: "test", Description: "test", Status: "test"},
			wantErr: false,
		},
		{
			name:    "Invalid Task: Error is expected",
			args:    models.Task{ID: 0, Title: "test", Description: "test", Status: "test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockStorage.Set(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("taskStorage.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
