package quizzer

import "time"

type Question struct {
	ID             string
	QuizID         string
	Question       string
	Index          int
	TimeLimit      time.Duration
	Answers        []string
	CorrectAnswer  string
	VideoURL       *string
	VideoStartTime *time.Duration
	VideoEndTime   *time.Duration
}
