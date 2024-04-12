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
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice/dto"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/transport/http/httphandler"
	"github.com/yigithankarabulut/ConcurrentTaskService/pkg/util"
)

func TestUpdateInvalidMethod(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodPost, "/update", nil)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code, want %v got %v", http.StatusMethodNotAllowed, w.Code)
	}
	shouldContain := "method POST not allowed"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestUpdateWithTimeout(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithContextTimeout(time.Second*-1),
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: context.DeadlineExceeded,
		}),
	)
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"}`
	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want %v got %v", http.StatusGatewayTimeout, w.Code)
	}
	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestUpdateQueryParamNotRequired(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodPut, "/update?status=active", nil)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "query parameters not required"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestUpdateEmptyBody(t *testing.T) {
	handler := httphandler.New()
	req := httptest.NewRequest(http.MethodPut, "/update", nil)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "request body is empty"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestUpdateInvalidBody(t *testing.T) {
	handler := httphandler.New()
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"`
	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want %v got %v", http.StatusBadRequest, w.Code)
	}
	shouldContain := "unexpected EOF"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want %v got %v", shouldContain, w.Body.String())
	}
}

func TestUpdateErrUnknown(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrUnknown,
		}),
	)
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"}`
	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Update(w, req)

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

func TestUpdateErrIDNotFound(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrIDNotFound,
		}),
	)
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"}`
	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want %v got %v", http.StatusNotFound, w.Code)
	}
	shouldContain, err := json.Marshal(util.BasicError(customerror.ErrIDNotFound.Error(), http.StatusNotFound))
	if err != nil {
		t.Errorf("error while marshalling error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}

func TestUpdateErrUpdate(t *testing.T) {
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			submitErr: customerror.ErrUpdate,
		}),
	)
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"}`
	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want %v got %v", http.StatusNotFound, w.Code)
	}
	shouldContain, err := json.Marshal(util.BasicError(customerror.ErrUpdate.Error(), http.StatusNotFound))
	if err != nil {
		t.Errorf("error while marshalling error: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}

func TestUpdateSuccess(t *testing.T) {
	resp := util.ResponseData{
		Data: dto.TaskResponse{
			ID:          1,
			Title:       "title",
			Description: "description",
			Status:      "active",
		},
	}
	handler := httphandler.New(
		httphandler.WithLogger(logger),
		httphandler.WithPool(&mockTaskWorker{
			response: resp,
		}),
	)
	body := `{"id":1,"status":"active","name":"test","description":"test","title":"test"}`
	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("wrong status code, want %v got %v", http.StatusOK, w.Code)
	}
	shouldContain, err := json.Marshal(util.Response(http.StatusOK, resp))
	if err != nil {
		t.Errorf("error while marshalling response: %v", err)
	}
	if !strings.Contains(w.Body.String(), string(shouldContain)) {
		t.Errorf("wrong body message, want %v got %v", string(shouldContain), w.Body.String())
	}
}
