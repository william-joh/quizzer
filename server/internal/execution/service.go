package execution

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/postgres"
)

type Service interface {
	CreateExecution(ctx context.Context, quizId string, hostId string) (string, error)
	GetExecution(ctx context.Context, code string) (*Execution, error)

	Run()
	Stop()
}

func NewInMemory(db postgres.Database) Service {
	return &inMemoryService{
		db:         db,
		executions: map[string]*Execution{},
		done:       make(chan bool),
	}
}

type inMemoryService struct {
	db         postgres.Database
	executions map[string]*Execution
	done       chan bool
}

// Run is a method that should periodically check if there are any executions that are done and if so, clean them up.
func (s *inMemoryService) Run() {
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()

		for {
			select {
			case <-s.done:
				log.Debug().Msg("Stopping execution service")
				return
			case <-ticker.C:
				log.Trace().Msg("Checking for done executions")
				for code, execution := range s.executions {
					if execution.IsDone {
						log.Debug().Str("code", code).Msg("Execution is done, cleaning up")
						execution.Close()
						delete(s.executions, code)
					}

					// Check if the execution was created more than 1 hour ago and if so, clean it up
					if time.Since(execution.CreatedAt) > time.Hour {
						log.Debug().Str("code", code).Msg("Execution is older than 1 hour, cleaning up")
						execution.Close()
						delete(s.executions, code)
					}
				}
			}
		}
	}()
}

func (s *inMemoryService) Stop() {
	s.done <- true
}

func (s *inMemoryService) CreateExecution(ctx context.Context, quizId string, hostId string) (string, error) {
	execution := Execution{
		Phase:           PhaseLobby,
		CurrentQuestion: 0,
		CreatedAt:       time.Now(),
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
	execution.Run()
	return execution.Code, nil
}

func (s *inMemoryService) GetExecution(ctx context.Context, code string) (*Execution, error) {
	execution, ok := s.executions[code]
	if !ok {
		return nil, errors.New("execution not found")
	}
	return execution, nil
}

func generateCode() string {
	// Generate a random number between 100000 and 999999
	num := rand.Intn(900000) + 100000
	return strconv.Itoa(num)
}
