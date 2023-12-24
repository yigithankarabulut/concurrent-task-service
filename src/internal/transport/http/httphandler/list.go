package httphandler

import (
	"context"
	"errors"
	"fmt"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/customerror"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/pkg/constant"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/pkg/util"
	"net/http"
	"strings"
)

// @Tags Task
// @Summary List Tasks by Status.
// @Description This endpoint is used for retrieving a list of tasks based on their status.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param 	status query string true "Task Status to retrieve" ExampleRequest
// @Success 200 {array} dto.TaskResponse "Success Response Body. List of tasks matching the specified status."
// @Failure 400 {object} util.ErrorResponse "Error Bad Request Response. Invalid request parameters."
// @Failure 404 {object} util.ErrorResponse "Error Not Found Response. No tasks found with the specified status."
// @Failure 500 {object} util.ErrorResponse "Error Internal Server. Server encountered an error."
// @Router /list [get]
func (h *httpHandler) List(w http.ResponseWriter, r *http.Request) {
	var (
		req models.TaskJobModel
	)
	ctx, cancel := context.WithTimeout(r.Context(), h.CancelTimeout)
	defer cancel()
	if r.Method != http.MethodGet {
		h.JSON(
			w,
			http.StatusMethodNotAllowed,
			fmt.Sprintf(constant.ErrMethodNotAllowed, r.Method),
		)
		return
	}
	if len(r.URL.Query()) == 0 {
		h.JSON(w,
			http.StatusBadRequest,
			util.BasicError("query parameters required", http.StatusBadRequest),
		)
		return
	}
	stat := r.URL.Query().Get("status")
	if stat == "" {
		h.JSON(w,
			http.StatusBadRequest,
			util.BasicError("invalid query parameters", http.StatusBadRequest),
		)
		return
	}
	invalidChars := []string{"<", ">", "!", "=", "(", ")", ";", ":", ",", "\"", "/", "\\", "?", "@", "[", "]", "{", "}", " ", "."}
	for _, c := range invalidChars {
		if strings.Contains(stat, c) || len(stat) > 255 {
			h.JSON(w,
				http.StatusBadRequest,
				util.BasicError("invalid query parameters", http.StatusBadRequest),
			)
			return
		}
	}
	req.Status = stat
	req.JOB = "LIST"
	req.Context = ctx
	// @Step: Submit to Pool
	res, err := h.pool.Submit(req)
	if err != nil {
		// @Step: Handle Errors
		if errors.Is(err, context.DeadlineExceeded) {
			h.JSON(w,
				http.StatusGatewayTimeout,
				util.BasicError(constant.ErrContextDeadline, http.StatusGatewayTimeout),
			)
			return
		}
		var cusErr *customerror.Error
		if errors.As(err, &cusErr) {
			clientMessage := cusErr.Message
			if cusErr.Data != nil {
				data, ok := cusErr.Data.(string)
				if ok {
					clientMessage = clientMessage + ", " + data
				}
			}

			if cusErr.Loggable {
				h.Logger.Error("httphandler List service.List", "err", clientMessage)
			}

			if cusErr == customerror.ErrGetAll {
				h.JSON(w,
					http.StatusNotFound,
					util.BasicError(clientMessage, http.StatusNotFound),
				)
				return
			}
		}
		h.JSON(w,
			http.StatusInternalServerError,
			util.BasicError(err.Error(), http.StatusInternalServerError),
		)
		return
	}
	// @Step: Return Success Response
	h.JSON(w,
		http.StatusOK,
		util.Response(http.StatusOK, res),
	)
}
