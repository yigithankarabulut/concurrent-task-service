package taskservice_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice/dto"
)

func TestListWithCancel(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := dto.ListTaskRequest{
		Status: "todo",
	}
	if _, err := taskService.List(ctx, req); !errors.Is(err, ctx.Err()) {
		t.Errorf("expected error: %v, got: %v", ctx.Err(), err)
	}
}

func TestListWithStorageError(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{
		listErr: errStorageList,
	}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.ListTaskRequest{
		Status: "todo",
	}
	if _, err := taskService.List(context.Background(), req); !errors.Is(err, errStorageList) {
		t.Errorf("expected error: %v, got: %v", errStorageList, err)
	}
}

func TestList(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.ListTaskRequest{
		Status: "todo",
	}
	if _, err := taskService.List(context.Background(), req); err != nil {
		t.Errorf("expected error: %v, got: %v", nil, err)
	}
}
