package stats_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/krxzy37/todoapp-go/internal/core/domain"
	core_http_server "github.com/krxzy37/todoapp-go/internal/core/transport/http/server"
)

type StatsHTTPHandler struct {
	statsService StatsService
}

type StatsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatsHTTPHandler(statsService StatsService) *StatsHTTPHandler {
	return &StatsHTTPHandler{
		statsService: statsService,
	}
}

func (h *StatsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStats,
		},
	}
}
