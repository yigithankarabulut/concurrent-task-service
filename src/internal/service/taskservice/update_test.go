package taskservice_test

import (
	"context"
	"errors"
	. "github.com/yigithankarabulut/ConcurrentTaskService/src/internal/service/taskservice"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/service/taskservice/dto"
	"testing"
)

func TestUpdateWithCancel(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := dto.UpdateTaskRequest{
		ID:          1,
		Title:       "title",
		Description: "description",
		Status:      "status",
	}
	if _, err := taskService.Update(ctx, req); !errors.Is(err, ctx.Err()) {
		t.Errorf("expected error: %v, got: %v", ctx.Err(), err)
	}
}

func TestUpdateWithGetError(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{
		getErr: errStorageGet,
	}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.UpdateTaskRequest{
		ID:          1,
		Title:       "title",
		Description: "description",
		Status:      "status",
	}
	if _, err := taskService.Update(context.Background(), req); !errors.Is(err, errStorageGet) {
		t.Errorf("expected error: %v, got: %v", errStorageGet, err)
	}
}

func TestUpdateWithStorageError(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{
		updateErr: errStorageUpdate,
	}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.UpdateTaskRequest{
		ID:          1,
		Title:       "title",
		Description: "description",
		Status:      "status",
	}
	if _, err := taskService.Update(context.Background(), req); !errors.Is(err, errStorageUpdate) {
		t.Errorf("expected error: %v, got: %v", errStorageUpdate, err)
	}
}

func TestUpdate(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.UpdateTaskRequest{
		ID:          1,
		Title:       "title",
		Description: "description",
		Status:      "status",
	}
	if _, err := taskService.Update(context.Background(), req); err != nil {
		t.Errorf("expected error: %v, got: %v", nil, err)
	}
}
