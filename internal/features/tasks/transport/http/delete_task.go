package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/krxzy37/todoapp-go/internal/core/logger"
	core_http_request "github.com/krxzy37/todoapp-go/internal/core/transport/http/request"
	core_http_response "github.com/krxzy37/todoapp-go/internal/core/transport/http/response"
)

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete an existing task in the system by its ID
// @Tags tasks
// @Param id path int true "ID of the task to delete"
// @Success 204 "Task successfully deleted"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponceHandler(log, rw)

	taskID, err := core_http_request.GetIntPathParam(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(

			err,
			"failed to get taskID path value",
		)

		return
	}

	if err := h.taskService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)
		return
	}

	responseHandler.NoContentResponse()

}
