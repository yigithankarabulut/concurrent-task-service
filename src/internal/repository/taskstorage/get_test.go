package taskstorage_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/yigithankarabulut/ConcurrentTaskService/src/internal/repository/taskstorage"
	"testing"
)

func Test_taskStorage_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockStorage := NewTaskStorage(WithTaskDB(db))
	mock.ExpectQuery("SELECT id, title, description, status FROM tasks WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status"}).
			AddRow(1, "title", "description", "status"))

	tests := []struct {
		name    string
		args    uint
		wantErr bool
	}{
		{
			name:    "Id is 1 and no error is expected",
			args:    1,
			wantErr: false,
		},
		{
			name:    "Id is 0 and error is expected",
			args:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mockStorage.Get(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("taskStorage.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
