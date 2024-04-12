package taskservice

import (
	"context"
	"fmt"
	"strconv"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/customerror"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/models"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice/dto"
)

func (s *taskService) Set(ctx context.Context, req dto.SetTaskRequest) (dto.TaskResponse, error) {
	select {
	case <-ctx.Done():
		return dto.TaskResponse{}, ctx.Err()
	default:
		if _, err := s.taskStorage.Get(req.ID); err == nil {
			_id := strconv.Itoa(int(req.ID))
			return dto.TaskResponse{}, fmt.Errorf("service.Set storage.Get: %w", customerror.ErrIDExists.AddData("'"+_id+"' already exists in the database."))
		}
		task := models.Task{
			ID:          req.ID,
			Title:       req.Title,
			Description: req.Description,
			Status:      req.Status,
		}
		if err := s.taskStorage.Set(task); err != nil {
			return dto.TaskResponse{}, fmt.Errorf("service.Set storage.Set: %w", err)
		}
		return dto.TaskResponse{
			ID:          req.ID,
			Title:       req.Title,
			Description: req.Description,
			Status:      req.Status,
		}, nil
	}
}
