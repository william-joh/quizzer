package quizzer

import "time"

type Quiz struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}
