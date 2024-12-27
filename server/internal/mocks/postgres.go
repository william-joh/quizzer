package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/william-joh/quizzer/server/internal/postgres"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

var _ postgres.Database = &Database{}

type Database struct {
	mock.Mock
}

func (m *Database) Close() error {
	m.Called()
	return nil
}

func (m *Database) InTx(ctx context.Context, f func(postgres.Session) error) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *Database) Do(ctx context.Context) postgres.Session {
	args := m.Called(ctx)
	return args.Get(0).(postgres.Session)
}

var _ postgres.Session = &Session{}

type Session struct {
	mock.Mock
}

func (m *Session) CreateQuestion(ctx context.Context, q quizzer.Question) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *Session) GetQuestion(ctx context.Context, id string) (quizzer.Question, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(quizzer.Question), args.Error(1)
}

func (m *Session) ListQuestions(ctx context.Context, quizID string) ([]quizzer.Question, error) {
	args := m.Called(ctx, quizID)
	return args.Get(0).([]quizzer.Question), args.Error(1)
}

func (m *Session) UpdateQuestion(ctx context.Context, q quizzer.Question) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *Session) DeleteQuestion(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *Session) CreateUser(ctx context.Context, id, username, password string) error {
	args := m.Called(ctx, id, username, password)
	return args.Error(0)
}

func (m *Session) GetUser(ctx context.Context, username string) (quizzer.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(quizzer.User), args.Error(1)
}

func (m *Session) GetUserByUsername(ctx context.Context, username string) (quizzer.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(quizzer.User), args.Error(1)
}

func (m *Session) DeleteUser(ctx context.Context, username string) error {
	args := m.Called(ctx, username)
	return args.Error(0)
}

func (m *Session) CreateQuiz(ctx context.Context, id, title, createdBy string) error {
	args := m.Called(ctx, id, title, createdBy)
	return args.Error(0)
}

func (m *Session) GetQuiz(ctx context.Context, id string) (quizzer.Quiz, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(quizzer.Quiz), args.Error(1)
}

func (m *Session) ListQuizzes(ctx context.Context) ([]quizzer.Quiz, error) {
	args := m.Called(ctx)
	return args.Get(0).([]quizzer.Quiz), args.Error(1)
}

func (m *Session) DeleteQuiz(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *Session) CreateAuthSession(ctx context.Context, userID, sessionID string) error {
	args := m.Called(ctx, userID, sessionID)
	return args.Error(0)
}

func (m *Session) GetAuthSession(ctx context.Context, sessionID string) (string, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(string), args.Error(1)
}
