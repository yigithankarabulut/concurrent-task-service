package taskservice

import (
	"context"
	"fmt"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/service/taskservice/dto"
)

func (s *taskService) Update(ctx context.Context, req dto.UpdateTaskRequest) (dto.TaskResponse, error) {
	select {
	case <-ctx.Done():
		return dto.TaskResponse{}, ctx.Err()
	default:
		if _, err := s.taskStorage.Get(req.ID); err != nil {
			return dto.TaskResponse{}, fmt.Errorf("service.Update storage.Get: %w", err)
		}
		task := models.Task{
			ID:          req.ID,
			Title:       req.Title,
			Description: req.Description,
			Status:      req.Status,
		}
		if err := s.taskStorage.Update(task); err != nil {
			return dto.TaskResponse{}, fmt.Errorf("service.Update storage.Update: %w", err)
		}
		return dto.TaskResponse{
			ID:          req.ID,
			Title:       req.Title,
			Description: req.Description,
			Status:      req.Status,
		}, nil
	}
}
