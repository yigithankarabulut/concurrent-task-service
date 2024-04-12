package taskstorage

import (
	"fmt"
	"strconv"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/customerror"
	. "github.com/yigithankarabulut/ConcurrentTaskService/internal/models"
)

func (s *taskStorage) Get(id uint) (Task, error) {
	task := Task{}
	err := s.db.QueryRow("SELECT id, title, description, status FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Title, &task.Description, &task.Status)
	_id := strconv.Itoa(int(id))
	if err != nil {
		return Task{}, fmt.Errorf("%w", customerror.ErrIDNotFound.AddData("'"+_id+"' does not exist in the database."))
	}
	return task, nil
}
