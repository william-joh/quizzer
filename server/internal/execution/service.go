package execution

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/william-joh/quizzer/server/internal/postgres"
)

type Service interface {
	CreateExecution(ctx context.Context, quizId string, hostId string) (string, error)
	JoinQuiz(ctx context.Context, executionId string, participantId string, name string, conn *websocket.Conn) error
	SetPhase(ctx context.Context, executionId string, phase Phase) error
	SubmitAnswer(ctx context.Context, executionId string, participantId string, questionID string, answer string) error
	GetQuizState(ctx context.Context, id string, participantId string) (QuizState, error)

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

func (s *inMemoryService) SetPhase(ctx context.Context, executionId string, phase Phase) error {
	execution, ok := s.executions[executionId]
	if !ok {
		return errors.New("execution not found")
	}

	execution.Phase = phase
	return nil
}

func (s *inMemoryService) JoinQuiz(ctx context.Context, executionId string, participantId string, name string, conn *websocket.Conn) error {
	execution, ok := s.executions[executionId]
	if !ok {
		return errors.New("execution not found")
	}
	if execution.Phase != PhaseLobby {
		return errors.New("execution is not in lobby phase")
	}
	if _, ok := s.getParticipant(execution, participantId); ok {
		return errors.New("participant already joined")
	}

	execution.Participants = append(execution.Participants, Participant{
		ID:   participantId,
		Name: name,
		Conn: conn,
	})

	return nil
}

func (s *inMemoryService) SubmitAnswer(ctx context.Context, executionId string, participantId string, questionID string, answer string) error {
	execution, ok := s.executions[executionId]
	if !ok {
		return errors.New("execution not found")
	}
	if execution.Phase != PhaseQuestion {
		return errors.New("execution is not in question phase")
	}

	participant, ok := s.getParticipant(execution, participantId)
	if !ok {
		return errors.New("participant not joined")
	}

	participant.Answers[questionID] = answer
	return nil
}

func (s *inMemoryService) GetQuizState(ctx context.Context, executionId string, participantId string) (QuizState, error) {
	execution, ok := s.executions[executionId]
	if !ok {
		return QuizState{}, errors.New("execution not found")
	}

	quizState := QuizState{Phase: execution.Phase}

	if execution.Host.ID == participantId {
		payload, err := s.getHostPayload(execution)
		if err != nil {
			return QuizState{}, err
		}

		quizState.Payload = payload
		return quizState, nil
	}

	payload, err := getParticipantPayload(execution)
	if err != nil {
		return QuizState{}, err
	}

	quizState.Payload = payload
	return quizState, nil
}

func (s *inMemoryService) getParticipant(execution *Execution, participantId string) (Participant, bool) {
	for _, p := range execution.Participants {
		if p.ID == participantId {
			return p, true
		}
	}
	return Participant{}, false
}

func (s *inMemoryService) getHostPayload(execution *Execution) (interface{}, error) {
	switch execution.Phase {
	case PhaseLobby:
		return getHostLobbyPayload(execution)
	case PhaseQuestion:
		return getHostQuestionPayload(execution)
	case PhaseResults:
		return getHostResultsPayload(execution)
	default:
		return nil, fmt.Errorf("unknown phase: %s", execution.Phase)
	}
}

func getHostLobbyPayload(execution *Execution) (interface{}, error) {
	payload := struct {
		Code             string   `json:"code"`
		QuizTitle        string   `json:"quizTitle"`
		HostName         string   `json:"hostName"`
		ParticipantNames []string `json:"participantNames"`
	}{
		Code:      execution.Code,
		QuizTitle: execution.Quiz.Title,
		HostName:  execution.Host.Username,
	}

	for _, p := range execution.Participants {
		payload.ParticipantNames = append(payload.ParticipantNames, p.Name)
	}

	return payload, nil
}

func getHostQuestionPayload(execution *Execution) (interface{}, error) {
	payload := struct {
		Options []string `json:"options"`
	}{}

	q := execution.Questions[execution.CurrentQuestion]
	payload.Options = q.Answers

	return payload, nil
}

func getHostResultsPayload(execution *Execution) (interface{}, error) {
	payload := struct {
		NrQuestions int `json:"nrQuestions"`
		Results     []struct {
			Name      string `json:"name"`
			NrCorrect int    `json:"nrCorrect"`
		} `json:"results"`
	}{
		NrQuestions: execution.CurrentQuestion,
	}

	for _, p := range execution.Participants {
		nrCorrect := 0
		for i, q := range execution.Questions {
			if i >= execution.CurrentQuestion {
				break
			}

			answer, ok := p.Answers[q.ID]
			if !ok {
				continue
			}
			if slices.Contains(q.CorrectAnswers, answer) {
				nrCorrect++
			}
		}

		payload.Results = append(payload.Results, struct {
			Name      string `json:"name"`
			NrCorrect int    `json:"nrCorrect"`
		}{
			Name:      p.Name,
			NrCorrect: nrCorrect,
		})
	}

	return payload, nil
}

func getParticipantPayload(execution *Execution) (interface{}, error) {
	switch execution.Phase {
	case PhaseLobby:
		return getParticipantLobbyPayload(execution)
	case PhaseQuestion:
		return getParticipantQuestionPayload(execution)
	case PhaseResults:
		return getParticipantResultsPayload(execution)
	default:
		return nil, fmt.Errorf("unknown phase: %s", execution.Phase)
	}
}

func getParticipantLobbyPayload(execution *Execution) (interface{}, error) {
	payload := struct {
		QuizTitle string `json:"quizTitle"`
		HostName  string `json:"hostName"`
	}{
		QuizTitle: execution.Quiz.Title,
		HostName:  execution.Host.Username,
	}

	return payload, nil
}

func getParticipantQuestionPayload(execution *Execution) (interface{}, error) {
	q := execution.Questions[execution.CurrentQuestion]

	payload := struct {
		Options []string `json:"options"`
	}{
		Options: q.Answers,
	}

	return payload, nil
}

func getParticipantResultsPayload(execution *Execution) (interface{}, error) {
	payload := struct {
		NrQuestions int `json:"nrQuestions"`
		NrCorrect   int `json:"nrCorrect"`
	}{
		NrQuestions: execution.CurrentQuestion,
	}

	return payload, nil
}
