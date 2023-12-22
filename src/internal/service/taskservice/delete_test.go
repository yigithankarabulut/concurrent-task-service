package taskservice_test

import (
	"context"
	"errors"
	. "github.com/yigithankarabulut/ConcurrentTaskService/src/internal/service/taskservice"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/service/taskservice/dto"
	"testing"
)

func TestDeleteWithCancel(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := dto.DeleteTaskRequest{
		ID: 1,
	}
	if err := taskService.Delete(ctx, req); !errors.Is(err, ctx.Err()) {
		t.Errorf("expected error: %v, got: %v", ctx.Err(), err)
	}
}

func TestDeleteWithGetError(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{
		getErr: errStorageGet,
	}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.DeleteTaskRequest{
		ID: 1,
	}
	if err := taskService.Delete(context.Background(), req); !errors.Is(err, errStorageGet) {
		t.Errorf("expected error: %v, got: %v", errStorageGet, err)
	}
}

func TestDeleteWithStorageError(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{
		deleteErr: errStorageDelete,
	}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.DeleteTaskRequest{
		ID: 1,
	}
	if err := taskService.Delete(context.Background(), req); !errors.Is(err, errStorageDelete) {
		t.Errorf("expected error: %v, got: %v", errStorageDelete, err)
	}
}

func TestDelete(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.DeleteTaskRequest{
		ID: 1,
	}
	if err := taskService.Delete(context.Background(), req); err != nil {
		t.Errorf("expected error: %v, got: %v", nil, err)
	}
}
