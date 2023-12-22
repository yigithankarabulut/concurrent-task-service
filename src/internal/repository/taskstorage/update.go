package taskstorage

import (
	"fmt"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/customerror"
	. "github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
	"strconv"
)

func (s *taskStorage) Update(task Task) error {
	_, err := s.db.Exec("UPDATE tasks SET title = ?, description = ?, status = ? WHERE id = ?", task.Title, task.Description, task.Status, task.ID)
	_id := strconv.Itoa(int(task.ID))
	if err != nil {
		return fmt.Errorf("%w", customerror.ErrUpdate.AddData("'"+_id+"' could not be updated."))
	}
	return nil
}
