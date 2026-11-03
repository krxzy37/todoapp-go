package stats_postgres_repository


type StatsRepository struct {

	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error) 

}