package stats_postgres_repository

import core_postgres_pool "github.com/krxzy37/todoapp-go/internal/core/repository/postgres/conn"

type StatsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatsRepository(
	pool core_postgres_pool.Pool,
) *StatsRepository {
	return &StatsRepository{
		pool: pool,
	}
}
