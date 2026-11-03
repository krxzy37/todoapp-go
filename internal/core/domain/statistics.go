package domain

import "time"

type Statistics struct {
	TasksCreated               int
	TaskCompleted              int
	TaskCompletedRate          *float64
	TasksAverageCompletionTIme *time.Duration
}

func NewStatistics(
	tasksCreated int,
	tasksCompleted int,
	tasksCompleteRate *float64,
	tasksAverageCompletionTIme *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:               tasksCreated,
		TaskCompleted:              tasksCompleted,
		TaskCompletedRate:          tasksCompleteRate,
		TasksAverageCompletionTIme: tasksAverageCompletionTIme,
	}
}
