package httphandler_test

import (
	"context"
	"encoding/json"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/customerror"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/transport/http/httphandler"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/pkg/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetInvalidMethod(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodPost, "/get", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code, want %v got %v", http.StatusMethodNotAllowed, w.Code)
	}
	shouldContain := "method POST not allowed"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestGetQueryParamRequired(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "query parameters required"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestGetInvalidQueryParam(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodGet, "/get?id=invalid", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "invalid query parameters"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestGetInvalidQueryParam2(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodGet, "/get?id=0", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "invalid query parameters"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestGetTimeout(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithContextTimeout(time.Second*-1),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: context.DeadlineExceeded,
		}),
	)
	req := httptest.NewRequest(http.MethodGet, "/get?id=1", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want %v got %v", http.StatusGatewayTimeout, w.Code)
	}
	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestGetErrUnknown(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithService(&mockTaskService{}),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrUnknown,
		}),
	)
	req := httptest.NewRequest(http.MethodGet, "/get?id=1", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code, want %v got %v", http.StatusInternalServerError, w.Code)
	}
	shouldContain, err := json.Marshal(util.BasicError(customerror.ErrUnknown.Error(), http.StatusInternalServerError))
	if err != nil {
		t.Errorf("error while response casting error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestGetErrIDNotFound(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithService(&mockTaskService{}),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrIDNotFound,
		}),
	)
	req := httptest.NewRequest(http.MethodGet, "/get?id=1", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want %v got %v", http.StatusNotFound, w.Code)
	}
	shouldContain, err := json.Marshal(util.BasicError(customerror.ErrIDNotFound.Error(), http.StatusNotFound))
	if err != nil {
		t.Errorf("error while response casting error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestGetSuccess(t *testing.T) {
	resp := util.ResponseData{
		Data: models.Task{
			ID:          1,
			Title:       "title",
			Description: "description",
			Status:      "status",
		},
		Status: http.StatusOK,
	}
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithService(&mockTaskService{}),
		httphandler.WithPool(&mockTaskWorker{
			response: resp,
		}),
	)
	req := httptest.NewRequest(http.MethodGet, "/get?id=1", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

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
