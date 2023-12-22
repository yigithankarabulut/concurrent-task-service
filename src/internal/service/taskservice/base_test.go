package taskservice_test

import (
	"errors"
	. "github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
)

var (
	errStorageDelete = errors.New("storage delete error")
	errStorageGet    = errors.New("storage get error")
	errStorageList   = errors.New("storage list error")
	errStorageSet    = errors.New("storage set error")
	errStorageUpdate = errors.New("storage update error")
)

type mockTaskStorage struct {
	deleteErr error
	getErr    error
	listErr   error
	setErr    error
	updateErr error
}

func (m *mockTaskStorage) Delete(uint) error {
	return m.deleteErr
}

func (m *mockTaskStorage) Get(uint) (Task, error) {
	return Task{}, m.getErr
}

func (m *mockTaskStorage) List(string) ([]Task, error) {
	return []Task{}, m.listErr
}

func (m *mockTaskStorage) Set(Task) error {
	return m.setErr
}

func (m *mockTaskStorage) Update(Task) error {
	return m.updateErr
}
