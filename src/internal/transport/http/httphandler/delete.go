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
// @Summary Delete Task by ID.
// @Description This endpoint is used for deleting a task based on its ID.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query integer true "Task ID required to delete"
// @Success 200 {object} string "Success Response Body Delete Successfully."
// @Failure 400 {object} util.ErrorResponse "Bad Request Response. Invalid request parameters."
// @Failure 404 {object} util.ErrorResponse "Not Found Response. No task found with the specified ID."
// @Failure 500 {object} util.ErrorResponse "Internal Server Error. Server encountered an error."
// @Router /delete [delete]
func (h *httpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		req models.TaskJobModel
	)
	ctx, cancel := context.WithTimeout(r.Context(), h.CancelTimeout)
	defer cancel()
	if r.Method != http.MethodDelete {
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
	req.JOB = "DELETE"
	req.Context = ctx

	// @Step: Submit to Pool
	if _, err = h.pool.Submit(req); err != nil {
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
				h.Logger.Error("httphandler Delete service.Delete", "err", clientMessage)
			}
			if cusErr == customerror.ErrIDNotFound {
				h.JSON(w,
					http.StatusNotFound,
					util.BasicError(clientMessage, http.StatusNotFound),
				)
				return
			}
			if cusErr == customerror.ErrDelete {
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
		util.Response(http.StatusOK, constant.DeletedSuccessfully),
	)
}
