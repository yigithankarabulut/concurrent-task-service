package taskstorage_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/models"
	. "github.com/yigithankarabulut/ConcurrentTaskService/internal/repository/taskstorage"
)

func Test_taskStorage_Update(t *testing.T) {
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
		mock.ExpectExec("UPDATE tasks").
			WithArgs(task.Title, task.Description, task.Status, task.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	tests := []struct {
		name    string
		args    models.Task
		wantErr bool
	}{
		{
			name:    "Task is valid and no error is expected",
			args:    models.Task{ID: 1, Title: "title", Description: "description", Status: "status"},
			wantErr: false,
		},
		{
			name:    "Task is valid and no error is expected",
			args:    models.Task{ID: 2, Title: "test", Description: "test", Status: "test"},
			wantErr: false,
		},
		{
			name:    "Task is invalid and error is expected",
			args:    models.Task{ID: 0, Title: "test", Description: "test", Status: "test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockStorage.Update(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("taskStorage.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
