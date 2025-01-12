package mocks

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
	"github.com/william-joh/quizzer/server/internal/execution"
)

var _ execution.Service = &ExecutionService{}

type ExecutionService struct {
	mock.Mock
}

func (m *ExecutionService) CreateExecution(ctx context.Context, quizId string, hostId string) (string, error) {
	args := m.Called(ctx, quizId, hostId)
	return args.String(0), args.Error(1)
}

func (m *ExecutionService) JoinQuiz(ctx context.Context, executionId string, participantId string, name string, conn *websocket.Conn) error {
	args := m.Called(ctx, executionId, participantId, name, conn)
	return args.Error(0)
}

func (m *ExecutionService) SetPhase(ctx context.Context, executionId string, phase execution.Phase) error {
	args := m.Called(ctx, executionId, phase)
	return args.Error(0)
}

func (m *ExecutionService) SubmitAnswer(ctx context.Context, executionId string, participantId string, questionID string, answer string) error {
	args := m.Called(ctx, executionId, participantId, questionID, answer)
	return args.Error(0)
}

func (m *ExecutionService) GetQuizState(ctx context.Context, id string, participantId string) (execution.QuizState, error) {
	args := m.Called(ctx, id, participantId)
	return args.Get(0).(execution.QuizState), args.Error(1)
}

func (m *ExecutionService) HandleMessages(ctx context.Context, conn *websocket.Conn) error {
	args := m.Called(ctx, conn)
	return args.Error(0)
}

func (m *ExecutionService) GetExecution(ctx context.Context, code string) (*execution.Execution, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*execution.Execution), args.Error(1)
}
