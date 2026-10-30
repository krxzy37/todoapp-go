package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/krxzy37/todoapp-go/internal/core/domain"
	core_logger "github.com/krxzy37/todoapp-go/internal/core/logger"
	core_http_request "github.com/krxzy37/todoapp-go/internal/core/transport/http/request"
	core_http_response "github.com/krxzy37/todoapp-go/internal/core/transport/http/response"
	core_http_types "github.com/krxzy37/todoapp-go/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`

	Completed core_http_types.Nullable[bool] `json:"completed"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`title` cant be null")
		}

		titleLen := len([]rune(*r.Title.Value))

		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`title` must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descLen := len([]rune(*r.Description.Value))

			if descLen > 1 || descLen > 1000 {
				return fmt.Errorf("`description` must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`completed` cant be null")
		}
	}

	return nil
}

type PatchUserResponse TaskDTOResponse

func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandker := core_http_response.NewHTTPResponceHandler(log, rw)

	taskID, err := core_http_request.GetIntPathParam(r, "id")
	if err != nil {

		responseHandker.ErrorResponse(
			err,
			"failed to get taskdID path value",
		)
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandker.ErrorResponse(err, "failed to decode and validate request")
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.taskService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandker.ErrorResponse(
			err,
			"failed to patch task",
		)
	}

	response := PatchUserResponse(taskDTOFromDomain(taskDomain))

	responseHandker.JSONResponse(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {

	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)

}
