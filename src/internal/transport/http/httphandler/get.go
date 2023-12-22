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
	"strconv"
)

// @Tags Task
// @Summary Get Task by ID.
// @Description This endpoint is used for retrieving a task based on its ID.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query integer true "Task ID to retrieve" ExampleRequest
// @Success 200 {object} dto.TaskResponse "Success Response Body. Task details with the specified ID."
// @Failure 400 {object} util.ErrorResponse "Error Bad Request Response. Invalid request parameters."
// @Failure 404 {object} util.ErrorResponse "Error Not Found Response. No task found with the specified ID."
// @Failure 500 {object} util.ErrorResponse "Error Internal Server. Server encountered an error."
// @Router /get [get]
func (h *httpHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	// @Step: Check Query Params
	if len(r.URL.Query()) == 0 {
		h.JSON(w,
			http.StatusBadRequest,
			util.BasicError("query parameters required", http.StatusBadRequest),
		)
		return
	}
	_id := r.URL.Query().Get("id")
	id, err := strconv.Atoi(_id)
	if err != nil {
		h.JSON(w,
			http.StatusBadRequest,
			util.BasicError("invalid query parameters", http.StatusBadRequest),
		)
		return
	}
	if id == 0 {
		h.JSON(w,
			http.StatusBadRequest,
			util.BasicError("invalid query parameters", http.StatusBadRequest),
		)
		return
	}

	req.ID = uint(id)
	req.JOB = "GET"
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
				h.Logger.Error("httphandler Get service.Get", "err", clientMessage)
			}
			if cusErr == customerror.ErrIDNotFound {
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
	return
}
