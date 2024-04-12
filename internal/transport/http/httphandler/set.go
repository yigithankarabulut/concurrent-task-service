package httphandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/yigithankarabulut/ConcurrentTaskService/internal/customerror"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/models"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/service/taskservice/dto"
	"github.com/yigithankarabulut/ConcurrentTaskService/internal/transport/http/basehttphandler"
	"github.com/yigithankarabulut/ConcurrentTaskService/pkg/constant"
	"github.com/yigithankarabulut/ConcurrentTaskService/pkg/util"
)

// @Tags Task
// @Summary 		Task Create.
// @Description 	This endpoint is used for creating a new task.
// @Accept			json
// @Produce			json
// @Security		BearerAuth
// @Param 			request body dto.SetTaskRequest true "Task Set Request Body"
// @Success 		200 {object} dto.TaskResponse "Success Response Body"
// @Failure 		400 {object} util.ErrorResponse "Error Bad Request Response"
// @Failure 		404 {object} util.ErrorResponse "Error Not Found Response"
// @Failure 		500 {object} util.ErrorResponse "Error Internal Server"
// @Router 			/set [post]
func (h *httpHandler) Set(w http.ResponseWriter, r *http.Request) {
	var (
		req models.TaskJobModel
	)
	ctx, cancel := context.WithTimeout(context.Background(), h.CancelTimeout)
	defer cancel()
	if r.Method != http.MethodPost {
		h.JSON(
			w,
			http.StatusMethodNotAllowed,
			fmt.Sprintf(constant.ErrMethodNotAllowed, r.Method),
		)
		return
	}
	if len(r.URL.Query()) > 0 {
		h.JSON(w,
			http.StatusBadRequest,
			util.BasicError("query parameters not required", http.StatusBadRequest),
		)
		return
	}
	// @Step: Validate Request
	resp, err := basehttphandler.Validate[dto.SetTaskRequest](r)
	if err != nil {
		h.JSON(w,
			http.StatusBadRequest,
			util.BasicError(err.Error(), http.StatusBadRequest),
		)
		return
	}
	resp.(dto.SetTaskRequest).TaskJobMapper(&req)
	req.JOB = "SET"
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
				h.Logger.Error("httphandler Set service.Set", "err", clientMessage)
			}

			if cusErr == customerror.ErrIDExists {
				h.JSON(w,
					http.StatusNotFound,
					util.BasicError(clientMessage, http.StatusNotFound),
				)
				return
			}
			if cusErr == customerror.ErrSet {
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
