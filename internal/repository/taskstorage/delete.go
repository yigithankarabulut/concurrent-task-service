package taskstorage

import (
	"fmt"
	"strconv"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/customerror"
)

func (s *taskStorage) Delete(id uint) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	_id := strconv.Itoa(int(id))
	if err != nil {
		return fmt.Errorf("%w", customerror.ErrDelete.AddData("'"+_id+"' could not be deleted."))
	}
	return err
}
