package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/krxzy37/todoapp-go/internal/core/logger"
	core_http_request "github.com/krxzy37/todoapp-go/internal/core/transport/http/request"
	core_http_response "github.com/krxzy37/todoapp-go/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary List of tasks
// @Description View a list of tasks with optional pagination and/or filtering by task author ID
// @Tags tasks
// @Produce json
// @Param user_id query int false "Filter tasks by author ID"
// @Param limit query int false "Page size for tasks"
// @Param offset query int false "Page offset for tasks"
// @Success 200 {object} GetTasksResponse "List of tasks"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponceHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQuerryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/Limit/Offset query params",
		)
		return
	}

	tasksDomains, err := h.taskService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	response := GetTasksResponse(tasksDTOsFromDomains(tasksDomains))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func getUserIDLimitOffsetQuerryParams(r *http.Request) (*int, *int, *int, error) {

	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `user_id` query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `limit` querry param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `offset` querry param: %w", err)
	}

	return userID, limit, offset, nil

}
