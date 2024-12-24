package quizzer

import "time"

type Quiz struct {
	ID        string
	Title     string
	CreatedBy string
	CreatedAt time.Time
}
