package workerservice

import (
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/service/taskservice"
	"log/slog"
	"os"
	"sync"
)

type TaskWorker interface {
	Submit(models.TaskJobModel) (any, error)
}

type taskWorker struct {
	workerCount int
	logger      *slog.Logger
	service     taskservice.TaskService
	ReqChan     chan models.TaskJobModel
	ResChan     chan any
	ErrChan     chan error
	done        chan struct{}
	Wg          *sync.WaitGroup
	mu          sync.Mutex
}

type TaskWorkerOption func(*taskWorker)

func WithWorkerCount(workerCount int) TaskWorkerOption {
	return func(t *taskWorker) {
		t.workerCount = workerCount
	}
}

func WithWaitGroup(wg *sync.WaitGroup) TaskWorkerOption {
	return func(t *taskWorker) {
		t.Wg = wg
	}
}

func WithService(service taskservice.TaskService) TaskWorkerOption {
	return func(t *taskWorker) {
		t.service = service
	}
}

func WithChannel(reqChan chan models.TaskJobModel, resChan chan any, errChan chan error, done chan struct{}) TaskWorkerOption {
	return func(t *taskWorker) {
		t.ReqChan = reqChan
		t.ResChan = resChan
		t.ErrChan = errChan
		t.done = done
	}
}

func StartTaskWorker(opts ...TaskWorkerOption) TaskWorker {
	tw := &taskWorker{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
		mu:     sync.Mutex{},
	}
	for _, opt := range opts {
		opt(tw)
	}
	for i := 0; i < tw.workerCount; i++ {
		tw.Wg.Add(1)
		go tw.worker()
	}
	return tw
}
