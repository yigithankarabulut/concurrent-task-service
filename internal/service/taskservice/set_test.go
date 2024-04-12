package taskservice_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice/dto"
)

func TestSetWithCancel(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := dto.SetTaskRequest{
		ID:          1,
		Title:       "title",
		Description: "description",
		Status:      "todo",
	}
	if _, err := taskService.Set(ctx, req); !errors.Is(err, ctx.Err()) {
		t.Errorf("expected error: %v, got: %v", ctx.Err(), err)
	}
}

func TestSetWithGetError(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{
		setErr: errStorageGet,
		getErr: errStorageGet,
	}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.SetTaskRequest{
		ID:          1,
		Title:       "title",
		Description: "description",
		Status:      "todo",
	}
	if _, err := taskService.Set(context.Background(), req); !errors.Is(err, errStorageGet) {
		t.Errorf("expected error: %v, got: %v", errStorageGet, err)
	}
}

func TestSetWithStorageError(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{
		setErr: errStorageSet,
		getErr: errStorageSet,
	}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.SetTaskRequest{
		ID:          1,
		Title:       "title",
		Description: "description",
		Status:      "todo",
	}
	if _, err := taskService.Set(context.Background(), req); !errors.Is(err, errStorageSet) {
		t.Errorf("expected error: %v, got: %v", errStorageSet, err)
	}
}

func TestSet(t *testing.T) {
	mockTaskStorage := &mockTaskStorage{
		getErr: errStorageSet,
	}
	taskService := NewTaskService(WithTaskStorage(mockTaskStorage))

	req := dto.SetTaskRequest{
		ID:          1,
		Title:       "title",
		Description: "description",
		Status:      "todo",
	}
	if _, err := taskService.Set(context.Background(), req); err != nil {
		t.Errorf("expected error: %v, got: %v", nil, err)
	}
}
