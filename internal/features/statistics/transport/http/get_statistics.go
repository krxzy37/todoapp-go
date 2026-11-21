package stats_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/krxzy37/todoapp-go/internal/core/domain"
	core_logger "github.com/krxzy37/todoapp-go/internal/core/logger"
	core_http_request "github.com/krxzy37/todoapp-go/internal/core/transport/http/request"
	core_http_response "github.com/krxzy37/todoapp-go/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TaskCompleted              int      `json:"tasks_completed"`
	TaskCompletedRate          *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

// GetStatistics godoc
// @Summary Get statistics
// @Description Get task statistics with optional filtering by user_id and/or time period
// @Tags statistics
// @Produce json
// @Param user_id query int false "Filter statistics by specific user"
// @Param from query string false "Start of the statistics review period (inclusive), format: YYYY-MM-DD"
// @Param to query string false "End of the statistics review period (exclusive), format: YYYY-MM-DD"
// @Success 200 {object} GetStatisticsResponse "Successfully retrieved statistics"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /statistics [get]
func (h *StatsHTTPHandler) GetStats(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponceHandler(log, rw)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to read query params",
		)
		return
	}

	stats, err := h.statsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)
	}

	response := toDTOFromDomain(stats)

	responseHandler.JSONResponse(response, http.StatusOK)

}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTIme != nil {
		duration := statistics.TasksAverageCompletionTIme.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TaskCompleted:              statistics.TaskCompleted,
		TaskCompletedRate:          statistics.TaskCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {

	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}
	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}
	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
