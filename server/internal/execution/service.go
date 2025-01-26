package execution

import (
	"context"
	"errors"
	"math/rand"
	"strconv"

	"github.com/william-joh/quizzer/server/internal/postgres"
)

type Service interface {
	CreateExecution(ctx context.Context, quizId string, hostId string) (string, error)
	GetExecution(ctx context.Context, code string) (*Execution, error)
}

func NewInMemory(db postgres.Database) Service {
	return &inMemoryService{db: db, executions: map[string]*Execution{}}
}

type inMemoryService struct {
	db         postgres.Database
	executions map[string]*Execution
}

func generateCode() string {
	// Generate a random number between 100000 and 999999
	num := rand.Intn(900000) + 100000
	return strconv.Itoa(num)
}

func (s *inMemoryService) CreateExecution(ctx context.Context, quizId string, hostId string) (string, error) {
	execution := Execution{
		Phase:           PhaseLobby,
		CurrentQuestion: 0,
	}
	for i := range 100 {
		code := generateCode()
		if _, ok := s.executions[code]; !ok {
			execution.Code = code
			break
		}

		if i == 99 {
			return "", errors.New("failed to generate a unique code")
		}
	}

	err := s.db.InTx(ctx, func(s postgres.Session) error {
		quiz, err := s.GetQuiz(ctx, quizId)
		if err != nil {
			return err
		}
		execution.Quiz = quiz

		questions, err := s.ListQuestions(ctx, quizId)
		if err != nil {
			return err
		}
		execution.Questions = questions

		host, err := s.GetUser(ctx, hostId)
		if err != nil {
			return err
		}
		execution.Host = host

		return nil
	})
	if err != nil {
		return "", err
	}

	s.executions[execution.Code] = &execution
	return execution.Code, nil
}

func (s *inMemoryService) GetExecution(ctx context.Context, code string) (*Execution, error) {
	execution, ok := s.executions[code]
	if !ok {
		return nil, errors.New("execution not found")
	}
	return execution, nil
}
