package taskservice_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice/dto"
)

func TestGetWithCancel(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := dto.GetTaskRequest{
		ID: 1,
	}
	if _, err := taskService.Get(ctx, req); !errors.Is(err, ctx.Err()) {
		t.Errorf("expected error: %v, got: %v", ctx.Err(), err)
	}
}

func TestGetWithStorageError(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{
		getErr: errStorageGet,
	}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.GetTaskRequest{
		ID: 1,
	}
	if _, err := taskService.Get(context.Background(), req); !errors.Is(err, errStorageGet) {
		t.Errorf("expected error: %v, got: %v", errStorageGet, err)
	}
}

func TestGet(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.GetTaskRequest{
		ID: 1,
	}
	if _, err := taskService.Get(context.Background(), req); err != nil {
		t.Errorf("expected error: %v, got: %v", nil, err)
	}
}
