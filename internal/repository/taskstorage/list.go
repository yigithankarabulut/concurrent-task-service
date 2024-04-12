package taskstorage

import (
	"fmt"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/customerror"
	. "github.com/yigithankarabulut/ConcurrentTaskService/internal/models"
)

func (s *taskStorage) List(status string) ([]Task, error) {
	tasks := make([]Task, 0)
	rows, err := s.db.Query("SELECT id, title, description, status FROM tasks WHERE status = ?", status)
	if err != nil {
		return nil, fmt.Errorf("%w", customerror.ErrGetAll.AddData("'"+status+"' could not be listed."))
	}
	defer rows.Close()
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
		if err != nil {
			return nil, fmt.Errorf("%w", customerror.ErrGetAll.AddData("'"+status+"' could not be listed."))
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
