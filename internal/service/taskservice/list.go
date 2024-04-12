package taskservice

import (
	"context"
	"fmt"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice/dto"
)

func (s *taskService) List(ctx context.Context, req dto.ListTaskRequest) ([]dto.TaskResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		tasks, err := s.taskStorage.List(req.Status)
		if err != nil {
			return nil, fmt.Errorf("service.List storage.List: %w", err)
		}
		var taskResponses []dto.TaskResponse
		for _, task := range tasks {
			taskResponses = append(taskResponses, dto.TaskResponse{
				ID:          task.ID,
				Title:       task.Title,
				Description: task.Description,
				Status:      task.Status,
			})
		}
		return taskResponses, nil
	}
}
