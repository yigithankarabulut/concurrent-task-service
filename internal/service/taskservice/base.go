package taskservice

import (
	"context"

	. "github.com/yigithankarabulut/ConcurrentTaskService/internal/repository/taskstorage"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice/dto"
)

type TaskService interface {
	Set(context.Context, dto.SetTaskRequest) (dto.TaskResponse, error)
	Get(context.Context, dto.GetTaskRequest) (dto.TaskResponse, error)
	List(context.Context, dto.ListTaskRequest) ([]dto.TaskResponse, error)
	Update(context.Context, dto.UpdateTaskRequest) (dto.TaskResponse, error)
	Delete(context.Context, dto.DeleteTaskRequest) error
}

type taskService struct {
	taskStorage TaskStorer
}

type TaskServiceOption func(*taskService)

func WithTaskStorage(taskStorage TaskStorer) TaskServiceOption {
	return func(s *taskService) {
		s.taskStorage = taskStorage
	}
}

func NewTaskService(opts ...TaskServiceOption) TaskService {
	s := &taskService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
