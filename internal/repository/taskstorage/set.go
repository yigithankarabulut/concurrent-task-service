package taskstorage

import (
	"fmt"
	"strconv"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/customerror"
	. "github.com/yigithankarabulut/ConcurrentTaskService/internal/models"
)

func (s *taskStorage) Set(task Task) error {
	_, err := s.db.Exec("INSERT INTO tasks (id, title, description, status) VALUES (?, ?, ?, ?)", task.ID, task.Title, task.Description, task.Status)
	_id := strconv.Itoa(int(task.ID))
	if err != nil {
		return fmt.Errorf("%w", customerror.ErrSet.AddData("'"+_id+"' could not be set."))
	}
	return nil
}
