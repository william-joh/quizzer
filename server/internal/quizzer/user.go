package quizzer

import "time"

type User struct {
	ID         string
	Username   string
	Password   string
	SignupDate time.Time
}
