package taskservice

import (
	"context"
	"fmt"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/service/taskservice/dto"
)

func (s *taskService) Get(ctx context.Context, req dto.GetTaskRequest) (dto.TaskResponse, error) {
	select {
	case <-ctx.Done():
		return dto.TaskResponse{}, ctx.Err()
	default:
		task, err := s.taskStorage.Get(req.ID)
		if err != nil {
			return dto.TaskResponse{}, fmt.Errorf("service.Get storage.Get: %w", err)
		}
		return dto.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
		}, nil
	}
}
