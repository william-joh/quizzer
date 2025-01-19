package execution

import (
	"fmt"
	"slices"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

type Participant struct {
	Conn    *websocket.Conn
	ID      string            `json:"userId"`
	Name    string            `json:"name"`
	Answers map[string]string `json:"answers"`
}

type Phase string

const (
	PhaseLobby    Phase = "lobby"
	PhaseQuestion Phase = "question"
	PhaseResults  Phase = "results"
)

type Execution struct {
	Code            string             `json:"id"`
	Quiz            quizzer.Quiz       `json:"quiz"`
	Questions       []quizzer.Question `json:"questions"`
	Host            quizzer.User       `json:"host"`
	HostConn        *websocket.Conn
	Participants    []Participant `json:"participants"`
	Phase           Phase         `json:"phase"`
	CurrentQuestion int           `json:"currentQuestion"`
}

type QuizState struct {
	Phase   Phase       `json:"phase"`
	Payload interface{} `json:"payload"`
}

type Message struct {
	Type   string      `json:"type"`
	Code   string      `json:"code"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func (e *Execution) HandleMessages(conn *websocket.Conn) error {
	var msg Message
	if err := conn.ReadJSON(&msg); err != nil {
		log.Error().Err(err).Msg("Failed to read message")
		return fmt.Errorf("read message: %w", err)
	}

	fmt.Println("Received message:", msg)

	var err error
	switch msg.Type {
	case "Join":
		err = e.handleJoinMsg(conn, msg)
	case "Start":
		err = e.handleStartMsg(conn)
	case "FinishQuestion":
		err = e.handleFinishQuestionMsg(conn)
	case "NextQuestion":
		err = e.handleNextQuestionMsg(conn)
	case "AnswerQuestion":
		err = e.handleAnswerQuestionMsg(conn, msg)
	default:
		log.Error().Str("type", msg.Type).Msg("Unknown message type")
		err = fmt.Errorf("unknown message type: %s", msg.Type)
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to handle message")
		return fmt.Errorf("handle message: %w", err)
	}

	return nil
}

func (e *Execution) handleJoinMsg(conn *websocket.Conn, msg Message) error {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		log.Error().Msgf("Failed to parse data, expected map[string]string, got %T", msg.Data)
		return fmt.Errorf("parse data: expected map[string]string, got %T", msg.Data)
	}

	username, ok := data["username"]
	if !ok {
		log.Error().Msg("Username not provided")
		return fmt.Errorf("username not provided")
	}

	participantId, ok := data["id"]
	if !ok {
		log.Error().Msg("Participant ID not provided")
		return fmt.Errorf("participant ID not provided")
	}

	if participantId == e.Host.ID {
		e.HostConn = conn
	} else {
		_, ok := e.getParticipant(participantId.(string))
		if !ok {
			participant := Participant{
				Conn:    conn,
				ID:      participantId.(string),
				Name:    username.(string),
				Answers: make(map[string]string),
			}
			e.Participants = append(e.Participants, participant)
		}

		log.Error().Msg("Participant already joined")
		// TODO: return fmt.Errorf("participant already joined")
	}

	// Broadcast the new quiz state
	if err := e.broadcastQuizState(); err != nil {
		log.Error().Err(err).Msg("Failed to broadcast quiz state")
		return fmt.Errorf("broadcast quiz state: %w", err)
	}

	return nil
}

func (e *Execution) handleStartMsg(conn *websocket.Conn) error {
	if e.HostConn != conn {
		log.Error().Msg("Only the host can start the quiz")
		return fmt.Errorf("only the host can start the quiz")
	}

	e.Phase = PhaseQuestion

	// Broadcast the new quiz state
	if err := e.broadcastQuizState(); err != nil {
		log.Error().Err(err).Msg("Failed to broadcast quiz state")
		return fmt.Errorf("broadcast quiz state: %w", err)
	}

	return nil
}

func (e *Execution) handleFinishQuestionMsg(conn *websocket.Conn) error {
	if e.HostConn != conn {
		log.Error().Msg("Only the host can finish a question")
		return fmt.Errorf("only the host can finish a question")
	}

	e.Phase = PhaseResults
	e.CurrentQuestion++

	// Broadcast the new quiz state
	if err := e.broadcastQuizState(); err != nil {
		log.Error().Err(err).Msg("Failed to broadcast quiz state")
		return fmt.Errorf("broadcast quiz state: %w", err)
	}

	return nil
}

func (e *Execution) handleNextQuestionMsg(conn *websocket.Conn) error {
	if e.HostConn != conn {
		log.Error().Msg("Only the host can move to the next question")
		return fmt.Errorf("only the host can move to the next question")
	}

	e.Phase = PhaseQuestion

	// Broadcast the new quiz state
	if err := e.broadcastQuizState(); err != nil {
		log.Error().Err(err).Msg("Failed to broadcast quiz state")
		return fmt.Errorf("broadcast quiz state: %w", err)
	}

	return nil
}

func (e *Execution) handleAnswerQuestionMsg(_ *websocket.Conn, msg Message) error {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		log.Error().Msgf("Failed to parse data, expected map[string]string, got %T", msg.Data)
		return fmt.Errorf("parse data: expected map[string]string, got %T", msg.Data)
	}

	participantId, ok := data["id"]
	if !ok {
		log.Error().Msg("Participant ID not provided")
		return fmt.Errorf("participant ID not provided")
	}

	participant, ok := e.getParticipant(participantId.(string))
	if !ok {
		log.Error().Msg("Participant not found")
		return fmt.Errorf("participant not found")
	}

	answer, ok := data["answer"]
	if !ok {
		log.Error().Msg("Answer not provided")
		return fmt.Errorf("answer not provided")
	}

	participant.Answers[e.Questions[e.CurrentQuestion].ID] = answer.(string)

	// Broadcast the new quiz state
	if err := e.broadcastQuizState(); err != nil {
		log.Error().Err(err).Msg("Failed to broadcast quiz state")
		return fmt.Errorf("broadcast quiz state: %w", err)
	}

	return nil
}

func (e *Execution) broadcastQuizState() error {
	// Send quiz state to host
	hostPayload, err := e.getHostPayload()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get host payload")
		return err
	}

	if err := e.HostConn.WriteJSON(hostPayload); err != nil {
		log.Error().Err(err).Msg("Failed to send quiz state to host")
		return err
	}

	for _, p := range e.Participants {
		participantPayload, err := e.getParticipantPayload()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get participant payload")
			return err
		}

		if err := p.Conn.WriteJSON(participantPayload); err != nil {
			log.Error().Err(err).Msg("Failed to send quiz state to participant")
			return err
		}
	}

	return nil
}

func (e *Execution) getHostPayload() (interface{}, error) {
	switch e.Phase {
	case PhaseLobby:
		return e.getHostLobbyPayload()
	case PhaseQuestion:
		return e.getHostQuestionPayload()
	case PhaseResults:
		return e.getHostResultsPayload()
	default:
		return nil, fmt.Errorf("unknown phase: %s", e.Phase)
	}
}

func (e *Execution) getHostLobbyPayload() (interface{}, error) {
	payload := struct {
		QuizTitle        string   `json:"quizTitle"`
		HostName         string   `json:"hostName"`
		IsHost           bool     `json:"isHost"`
		ParticipantNames []string `json:"participantNames"`
		Phase            string   `json:"phase"`
	}{
		QuizTitle:        e.Quiz.Title,
		HostName:         e.Host.Username,
		IsHost:           true,
		Phase:            string(e.Phase),
		ParticipantNames: []string{},
	}

	for _, p := range e.Participants {
		payload.ParticipantNames = append(payload.ParticipantNames, p.Name)
	}

	return payload, nil
}

func (e *Execution) getHostQuestionPayload() (interface{}, error) {
	payload := struct {
		Question string   `json:"question"`
		Options  []string `json:"options"`
		Phase    string   `json:"phase"`
	}{
		Phase: string(e.Phase),
	}

	q := e.Questions[e.CurrentQuestion]
	payload.Options = q.Answers
	payload.Question = q.Question

	return payload, nil
}

func (e *Execution) getHostResultsPayload() (interface{}, error) {
	payload := struct {
		Phase                string `json:"phase"`
		NrQuestionsCompleted int    `json:"nrQuestionsCompleted"`
		TotalQuestions       int    `json:"totalQuestions"`
		Results              []struct {
			Name      string `json:"name"`
			NrCorrect int    `json:"nrCorrect"`
		} `json:"results"`
	}{
		Phase:                string(e.Phase),
		NrQuestionsCompleted: e.CurrentQuestion,
		TotalQuestions:       len(e.Questions),
	}

	for _, p := range e.Participants {
		nrCorrect := 0
		for i, q := range e.Questions {
			if i >= e.CurrentQuestion {
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

func (e *Execution) getParticipantPayload() (interface{}, error) {
	switch e.Phase {
	case PhaseLobby:
		return e.getParticipantLobbyPayload()
	case PhaseQuestion:
		return e.getParticipantQuestionPayload()
	case PhaseResults:
		return e.getParticipantResultsPayload()
	default:
		return nil, fmt.Errorf("unknown phase: %s", e.Phase)
	}
}

func (e *Execution) getParticipantLobbyPayload() (interface{}, error) {
	return struct {
		QuizTitle string `json:"quizTitle"`
		HostName  string `json:"hostName"`
		IsHost    bool   `json:"isHost"`
		Phase     string `json:"phase"`
	}{
		QuizTitle: e.Quiz.Title,
		HostName:  e.Host.Username,
		IsHost:    false,
		Phase:     string(e.Phase),
	}, nil
}

func (e *Execution) getParticipantQuestionPayload() (interface{}, error) {
	payload := struct {
		Options []string `json:"options"`
		Phase   string   `json:"phase"`
	}{
		Phase: string(e.Phase),
	}

	q := e.Questions[e.CurrentQuestion]
	payload.Options = q.Answers

	return payload, nil
}

func (e *Execution) getParticipantResultsPayload() (interface{}, error) {
	payload := struct {
		Phase string `json:"phase"`
	}{
		Phase: string(e.Phase),
	}
	return payload, nil
}

func (e *Execution) getParticipant(participantId string) (*Participant, bool) {
	for _, p := range e.Participants {
		if p.ID == participantId {
			return &p, true
		}
	}

	return nil, false
}
