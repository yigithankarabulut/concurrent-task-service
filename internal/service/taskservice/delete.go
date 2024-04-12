package taskservice

import (
	"context"
	"fmt"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice/dto"
)

func (s *taskService) Delete(ctx context.Context, req dto.DeleteTaskRequest) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if _, err := s.taskStorage.Get(req.ID); err != nil {
			return fmt.Errorf("service.Delete storage.Get: %w", err)
		}
		if err := s.taskStorage.Delete(req.ID); err != nil {
			return fmt.Errorf("service.Delete storage.Delete: %w", err)
		}
		return nil
	}
}
