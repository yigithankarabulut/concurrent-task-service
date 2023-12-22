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

func TestSetInvalidMethod(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodGet, "/set", nil)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code, want %v got %v", http.StatusMethodNotAllowed, w.Code)
	}
	shouldContain := "method GET not allowed"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestSetWithTimeout(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithContextTimeout(time.Second*-1),
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: context.DeadlineExceeded,
		},
		))
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want %v got %v", http.StatusGatewayTimeout, w.Code)
	}
	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestSetQueryParamNotRequired(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodPost, "/set?status=active", nil)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "query parameters not required"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestSetInvalidBody(t *testing.T) {
	handler := httphandler.New()
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"`
	req := httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "unexpected EOF"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestSetEmptyBody(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodPost, "/set", nil)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "request body is empty"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestSetErrUnknown(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrUnknown,
		},
		))
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code, want %v got %v", http.StatusInternalServerError, w.Code)
	}
	shouldContain, err := json.Marshal(util.BasicError(customerror.ErrUnknown.Error(), http.StatusInternalServerError))
	if err != nil {
		t.Errorf("error while marshalling error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}

func TestSetErrIDExists(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrIDExists,
		},
		))
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want %v got %v", http.StatusNotFound, w.Code)
	}
	shouldContain, err := json.Marshal(util.BasicError(customerror.ErrIDExists.Error(), http.StatusNotFound))
	if err != nil {
		t.Errorf("error while marshalling error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}

func TestSetErrSet(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrSet,
		},
		))
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want %v got %v", http.StatusNotFound, w.Code)
	}
	shouldContain, err := json.Marshal(util.BasicError(customerror.ErrSet.Error(), http.StatusNotFound))
	if err != nil {
		t.Errorf("error while marshalling error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}

func TestSetSuccess(t *testing.T) {
	resp := util.ResponseData{
		Data: dto.TaskResponse{
			ID:          1,
			Status:      "active",
			Description: "test",
			Title:       "test",
		},
		Status: http.StatusOK,
	}
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			response: resp,
		},
		))
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("wrong status code, want %v got %v", http.StatusOK, w.Code)
	}

	shouldContain, err := json.Marshal(util.Response(http.StatusOK, resp))
	if err != nil {
		t.Errorf("error while marshalling error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}
