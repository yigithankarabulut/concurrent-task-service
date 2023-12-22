package taskstorage

import (
	"fmt"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/customerror"
	. "github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
	"strconv"
)

func (s *taskStorage) Set(task Task) error {
	_, err := s.db.Exec("INSERT INTO tasks (id, title, description, status) VALUES (?, ?, ?, ?)", task.ID, task.Title, task.Description, task.Status)
	_id := strconv.Itoa(int(task.ID))
	if err != nil {
		return fmt.Errorf("%w", customerror.ErrSet.AddData("'"+_id+"' could not be set."))
	}
	return nil
}
