package httphandler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/customerror"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/transport/http/httphandler"
	"github.com/yigithankarabulut/ConcurrentTaskService/pkg/constant"
	"github.com/yigithankarabulut/ConcurrentTaskService/pkg/util"
)

func TestDeleteInvalidMethod(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodGet, "/delete", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code, want %v got %v", http.StatusMethodNotAllowed, w.Code)
	}
	shouldContain := "method GET not allowed"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestDeleteWithTimeout(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithContextTimeout(time.Second*-1),
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: context.DeadlineExceeded,
		}),
	)
	req := httptest.NewRequest(http.MethodDelete, "/delete?id=1", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want %v got %v", http.StatusGatewayTimeout, w.Code)
	}
	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestDeleteQueryParamRequired(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "query parameters required"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestDeleteInvalidQueryParam(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodDelete, "/delete?id=invalid", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "invalid query parameters"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestDeleteInvalidQueryParam2(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodDelete, "/delete?id=0", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "invalid query parameters"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestDeleteTimeout(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithContextTimeout(time.Second*-1),
		httphandler.WithService(&mockTaskService{}),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: context.DeadlineExceeded,
		}),
	)
	req := httptest.NewRequest(http.MethodDelete, "/delete?id=1", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want %v got %v", http.StatusGatewayTimeout, w.Code)
	}
	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestDeleteErrUnknown(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithService(&mockTaskService{}),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrUnknown,
		}),
		httphandler.WithLogger(logger),
	)
	req := httptest.NewRequest(http.MethodDelete, "/delete?id=1", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

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

func TestDeleteErrIDNotFound(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithService(&mockTaskService{}),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrIDNotFound,
		}),
		httphandler.WithLogger(logger),
	)
	req := httptest.NewRequest(http.MethodDelete, "/delete?id=1", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

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

func TestDeleteSuccess(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithService(&mockTaskService{}),
		httphandler.WithPool(&mockTaskWorker{}),
		httphandler.WithLogger(logger),
	)
	req := httptest.NewRequest(http.MethodDelete, "/delete?id=1", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("wrong status code, want %v got %v", http.StatusOK, w.Code)
	}
	shouldContain, err := json.Marshal(util.Response(http.StatusOK, constant.DeletedSuccessfully))
	if err != nil {
		t.Errorf("error while response casting error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}
