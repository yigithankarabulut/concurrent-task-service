package httphandler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/workerservice"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/transport/http/basehttphandler"
)

type HTTPHandler interface {
	Set(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
}

type httpHandler struct {
	service taskservice.TaskService
	pool    workerservice.TaskWorker
	basehttphandler.Handler
}

type StoreHandlerOption func(*httpHandler)

func WithPool(pool workerservice.TaskWorker) StoreHandlerOption {
	return func(handler *httpHandler) {
		handler.pool = pool
	}
}

func WithService(service taskservice.TaskService) StoreHandlerOption {
	return func(handler *httpHandler) {
		handler.service = service
	}
}

func WithContextTimeout(d time.Duration) StoreHandlerOption {
	return func(handler *httpHandler) {
		handler.CancelTimeout = d
	}
}

func WithLogger(l *slog.Logger) StoreHandlerOption {
	return func(handler *httpHandler) {
		handler.Logger = l
	}
}

func New(opts ...StoreHandlerOption) HTTPHandler {
	handler := &httpHandler{
		Handler: basehttphandler.Handler{},
	}
	for _, opt := range opts {
		opt(handler)
	}
	return handler
}
