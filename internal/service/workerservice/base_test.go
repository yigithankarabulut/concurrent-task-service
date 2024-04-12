package workerservice_test

import (
	"context"
	"errors"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice/dto"
)

var (
	errServiceDelete = errors.New("service delete error")
	errServiceGet    = errors.New("service get error")
	errServiceList   = errors.New("service list error")
	errServiceSet    = errors.New("service set error")
	errServiceUpdate = errors.New("service update error")
)

type mockTaskService struct {
	deleteErr error
	getErr    error
	listErr   error
	setErr    error
	updateErr error
}

func (m *mockTaskService) Delete(context.Context, dto.DeleteTaskRequest) error {
	return m.deleteErr
}

func (m *mockTaskService) Get(context.Context, dto.GetTaskRequest) (dto.TaskResponse, error) {
	return dto.TaskResponse{}, m.getErr
}

func (m *mockTaskService) List(context.Context, dto.ListTaskRequest) ([]dto.TaskResponse, error) {
	return []dto.TaskResponse{}, m.listErr
}

func (m *mockTaskService) Set(context.Context, dto.SetTaskRequest) (dto.TaskResponse, error) {
	return dto.TaskResponse{}, m.setErr
}

func (m *mockTaskService) Update(context.Context, dto.UpdateTaskRequest) (dto.TaskResponse, error) {
	return dto.TaskResponse{}, m.updateErr
}
