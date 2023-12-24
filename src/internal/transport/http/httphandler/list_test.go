package httphandler_test

import (
	"context"
	"encoding/json"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/customerror"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/service/taskservice/dto"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/transport/http/httphandler"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/pkg/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestListInvalidMethod(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodDelete, "/list", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code, want %v got %v", http.StatusMethodNotAllowed, w.Code)
	}
	shouldContain := "method " + http.MethodDelete + " not allowed"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestListWithTimeout(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithContextTimeout(time.Second*-1),
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: context.DeadlineExceeded,
		}),
	)
	req := httptest.NewRequest(http.MethodGet, "/list?status=active", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want %v got %v", http.StatusGatewayTimeout, w.Code)
	}
	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestListQueryParamRequired(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "query parameters required"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestListInvalidParams(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodGet, "/list?status=./...", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain, err := json.Marshal(util.BasicError("invalid query parameters", http.StatusBadRequest))
	if err != nil {
		t.Errorf("error while marshalling error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}

func TestListErrUnknown(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrUnknown,
		}),
	)
	req := httptest.NewRequest(http.MethodGet, "/list?status=okey", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code, want %v got %v", http.StatusInternalServerError, w.Code)
	}
	shouldContain, err := json.Marshal(util.BasicError(customerror.ErrUnknown.Error(), http.StatusInternalServerError))
	if err != nil {
		t.Errorf("error while response casting error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}

func TestListErrGetAll(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrGetAll,
		}),
	)
	req := httptest.NewRequest(http.MethodGet, "/list?status=passive", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want %v got %v", http.StatusNotFound, w.Code)
	}
	shouldContain, err := json.Marshal(util.BasicError(customerror.ErrGetAll.Error(), http.StatusNotFound))
	if err != nil {
		t.Errorf("error while response casting error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}

func TestListSuccess(t *testing.T) {
	resp := util.ResponseData{
		Data: []dto.TaskResponse{
			{
				ID:          1,
				Title:       "title",
				Description: "description",
				Status:      "active",
			},
			{
				ID:          2,
				Title:       "title2",
				Description: "description2",
				Status:      "active",
			},
		},
		Status: http.StatusOK,
	}
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			response: resp,
		},
		),
	)

	req := httptest.NewRequest(http.MethodGet, "/list?status=success", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("wrong status code, want %v got %v", http.StatusOK, w.Code)
	}
	shouldContain, err := json.Marshal(util.Response(http.StatusOK, resp))
	if err != nil {
		t.Errorf("error while response casting error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}
