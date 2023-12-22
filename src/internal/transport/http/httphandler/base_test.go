package httphandler_test

import (
	"context"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/service/taskservice/dto"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/pkg/util"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type mockTaskService struct {
	baseRes   dto.TaskResponse
	listRes   []dto.TaskResponse
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
	return m.baseRes, m.getErr
}

func (m *mockTaskService) List(context.Context, dto.ListTaskRequest) ([]dto.TaskResponse, error) {
	return m.listRes, m.listErr
}

func (m *mockTaskService) Set(context.Context, dto.SetTaskRequest) (dto.TaskResponse, error) {
	return m.baseRes, m.setErr
}

func (m *mockTaskService) Update(context.Context, dto.UpdateTaskRequest) (dto.TaskResponse, error) {
	return m.baseRes, m.updateErr
}

type mockTaskWorker struct {
	submitErr error
	response  util.ResponseData
}

func (m *mockTaskWorker) Submit(models.TaskJobModel) (any, error) {
	return m.response, m.submitErr
}
