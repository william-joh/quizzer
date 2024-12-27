package quizzer

import "time"

type Quiz struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}
