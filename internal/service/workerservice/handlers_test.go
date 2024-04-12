package workerservice_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/models"
	. "github.com/yigithankarabulut/ConcurrentTaskService/internal/service/workerservice"
)

var WokerCount = 50

func TestTaskWorkerWithCancel(t *testing.T) {
	reqChan := make(chan models.TaskJobModel, WokerCount)
	resChan := make(chan any, WokerCount)
	errChan := make(chan error, WokerCount)
	doneCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	mockService := &mockTaskService{}
	worker := StartTaskWorker(
		WithWorkerCount(WokerCount),
		WithWaitGroup(wg),
		WithService(mockService),
		WithChannel(reqChan, resChan, errChan, doneCh),
	)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	job := models.TaskJobModel{
		ID:      1,
		Context: ctx,
		JOB:     "GET",
	}
	_, err := worker.Submit(job)
	if err != nil && !errors.Is(err, ctx.Err()) {
		t.Errorf("expected error: %v, got: %v", ctx.Err(), err)
	}
	close(doneCh)
}

func TestTaskWorkerWithGet(t *testing.T) {
	reqChan := make(chan models.TaskJobModel, WokerCount)
	resChan := make(chan any, WokerCount)
	errChan := make(chan error, WokerCount)
	doneCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	mockService := &mockTaskService{
		getErr: errServiceGet,
	}
	worker := StartTaskWorker(
		WithWorkerCount(WokerCount),
		WithWaitGroup(wg),
		WithService(mockService),
		WithChannel(reqChan, resChan, errChan, doneCh),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	job := models.TaskJobModel{
		ID:      1,
		Context: ctx,
		JOB:     "GET",
	}
	if _, err := worker.Submit(job); !errors.Is(err, errServiceGet) {
		t.Errorf("expected error: %v, got: %v", errServiceGet, err)
	}
	close(doneCh)
}

func TestTaskWorkerWithSet(t *testing.T) {
	reqChan := make(chan models.TaskJobModel, WokerCount)
	resChan := make(chan any, WokerCount)
	errChan := make(chan error, WokerCount)
	doneCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	mockService := &mockTaskService{
		setErr: errServiceSet,
	}
	worker := StartTaskWorker(
		WithWorkerCount(WokerCount),
		WithWaitGroup(wg),
		WithService(mockService),
		WithChannel(reqChan, resChan, errChan, doneCh),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	job := models.TaskJobModel{
		ID:      1,
		Context: ctx,
		JOB:     "SET",
	}
	if _, err := worker.Submit(job); !errors.Is(err, errServiceSet) {
		t.Errorf("expected error: %v, got: %v", errServiceSet, err)
	}
	close(doneCh)
}

func TestTaskWorkerWithDelete(t *testing.T) {
	reqChan := make(chan models.TaskJobModel, WokerCount)
	resChan := make(chan any, WokerCount)
	errChan := make(chan error, WokerCount)
	doneCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	mockService := &mockTaskService{
		deleteErr: errServiceDelete,
	}
	worker := StartTaskWorker(
		WithWorkerCount(WokerCount),
		WithWaitGroup(wg),
		WithService(mockService),
		WithChannel(reqChan, resChan, errChan, doneCh),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	job := models.TaskJobModel{
		ID:      1,
		Context: ctx,
		JOB:     "DELETE",
	}
	if _, err := worker.Submit(job); !errors.Is(err, errServiceDelete) {
		t.Errorf("expected error: %v, got: %v", errServiceDelete, err)
	}
	close(doneCh)
}

func TestTaskWorkerWithUpdate(t *testing.T) {
	reqChan := make(chan models.TaskJobModel, WokerCount)
	resChan := make(chan any, WokerCount)
	errChan := make(chan error, WokerCount)
	doneCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	mockService := &mockTaskService{
		updateErr: errServiceUpdate,
	}
	worker := StartTaskWorker(
		WithWorkerCount(WokerCount),
		WithWaitGroup(wg),
		WithService(mockService),
		WithChannel(reqChan, resChan, errChan, doneCh),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	job := models.TaskJobModel{
		ID:      1,
		Context: ctx,
		JOB:     "UPDATE",
	}
	if _, err := worker.Submit(job); !errors.Is(err, errServiceUpdate) {
		t.Errorf("expected error: %v, got: %v", errServiceUpdate, err)
	}
	close(doneCh)
}

func TestTaskWorkerWithList(t *testing.T) {
	reqChan := make(chan models.TaskJobModel, WokerCount)
	resChan := make(chan any, WokerCount)
	errChan := make(chan error, WokerCount)
	doneCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	mockService := &mockTaskService{
		listErr: errServiceList,
	}
	worker := StartTaskWorker(
		WithWorkerCount(WokerCount),
		WithWaitGroup(wg),
		WithService(mockService),
		WithChannel(reqChan, resChan, errChan, doneCh),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	job := models.TaskJobModel{
		ID:      1,
		Context: ctx,
		JOB:     "LIST",
	}
	if _, err := worker.Submit(job); !errors.Is(err, errServiceList) {
		t.Errorf("expected error: %v, got: %v", errServiceList, err)
	}
	close(doneCh)
}

func TestTaskWorkerWithInvalidCRUD(t *testing.T) {
	reqChan := make(chan models.TaskJobModel, WokerCount)
	resChan := make(chan any, WokerCount)
	errChan := make(chan error, WokerCount)
	doneCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	mockService := &mockTaskService{}
	worker := StartTaskWorker(
		WithWorkerCount(WokerCount),
		WithWaitGroup(wg),
		WithService(mockService),
		WithChannel(reqChan, resChan, errChan, doneCh),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 1/2*time.Second)
	defer cancel()
	job := models.TaskJobModel{
		ID:      1,
		Context: ctx,
		JOB:     "INVALID",
	}
	if _, err := worker.Submit(job); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected error: %v, got: %v", context.DeadlineExceeded, err)
	}
	close(doneCh)
}
