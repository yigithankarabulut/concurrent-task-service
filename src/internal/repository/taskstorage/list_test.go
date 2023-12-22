package taskstorage_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/yigithankarabulut/ConcurrentTaskService/src/internal/repository/taskstorage"
	"testing"
)

func Test_taskStorage_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockStorage := NewTaskStorage(WithTaskDB(db))
	mock.ExpectQuery("SELECT id, title, description, status FROM tasks WHERE status = ?").
		WithArgs("active").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status"}).
		AddRow(1, "title", "description", "status"))

	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{
			name:    "Status is 'active' and no error is expected",
			args:    "active",
			wantErr: false,
		},
		{
			name:    "Status is 'status' and error is expected",
			args:    "status",
			wantErr: true,
		},
		{
			name:    "Status is empty and error is expected",
			args:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockStorage.List(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("taskStorage.List() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
