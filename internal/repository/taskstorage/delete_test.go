package taskstorage_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/models"
	. "github.com/yigithankarabulut/ConcurrentTaskService/internal/repository/taskstorage"
)

func Test_taskStorage_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockStorage := NewTaskStorage(WithTaskDB(db))

	mock.ExpectExec("DELETE FROM tasks WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	task := models.Task{
		ID:          1,
		Title:       "title",
		Description: "description",
		Status:      "status",
	}
	tests := []struct {
		name    string
		args    models.Task
		wantErr bool
	}{
		{
			name:    "Task is valid and no error is expected",
			args:    task,
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
			if err := mockStorage.Delete(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("taskStorage.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
