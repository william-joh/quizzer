package api_test

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func TestQuiz_CreateQuiz(t *testing.T) {
	client, sessionMock, server := authenticatedClient(t)
	defer server.Close()

	var quizID string
	quiz := `{
		"title": "testquiz",
		"questions": [
			{
				"question": "What is the capital of Sweden?",
				"answers": ["Stockholm", "Oslo", "Copenhagen", "Helsinki"],
				"correctAnswers": ["Stockholm"],
				"timeLimitSeconds": 20,
				"index": 0
			},
			{
				"question": "What is the capital of Norway?",
				"answers": ["Stockholm", "Oslo", "Copenhagen", "Helsinki"],
				"correctAnswers": ["Oslo"],
				"timeLimitSeconds": 10,
				"index": 1
			}
		]
	}`
	sessionMock.On("CreateQuiz", anyCtx, anyString, "testquiz", anyString).
		Return(nil)

	expectedQuestions := []quizzer.Question{
		{
			Question:         "What is the capital of Sweden?",
			Answers:          []string{"Stockholm", "Oslo", "Copenhagen", "Helsinki"},
			CorrectAnswers:   []string{"Stockholm"},
			Index:            0,
			TimeLimitSeconds: 20,
		},
		{
			Question:         "What is the capital of Norway?",
			Answers:          []string{"Stockholm", "Oslo", "Copenhagen", "Helsinki"},
			CorrectAnswers:   []string{"Oslo"},
			Index:            1,
			TimeLimitSeconds: 10,
		},
	}

	for i := range expectedQuestions {
		sessionMock.On("CreateQuestion", anyCtx, mock.MatchedBy(func(q quizzer.Question) bool {
			return q.Question == expectedQuestions[i].Question &&
				q.TimeLimitSeconds == expectedQuestions[i].TimeLimitSeconds &&
				slices.Equal(q.CorrectAnswers, expectedQuestions[i].CorrectAnswers) &&
				slices.Equal(q.Answers, expectedQuestions[i].Answers)
		})).
			Return(nil)
	}

	req, err := http.NewRequest("POST", server.URL+"/quizzes", strings.NewReader(quiz))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Get quizId from body
	var respBody struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	quizID = respBody.ID
	require.NotEmpty(t, quizID)

	mock.AssertExpectationsForObjects(t, sessionMock)
}

func TestQuiz_GetQuiz(t *testing.T) {
	client, sessionMock, server := authenticatedClient(t)
	defer server.Close()

	quizID := "quizID"
	sessionMock.On("GetQuiz", anyCtx, quizID).
		Return(quizzer.Quiz{ID: quizID, Title: "testquiz"}, nil)

	sessionMock.On("ListQuestions", anyCtx, quizID).Return([]quizzer.Question{
		{
			ID:               "questionID1",
			QuizID:           quizID,
			Question:         "What is the capital of Sweden?",
			Answers:          []string{"Stockholm", "Oslo", "Copenhagen", "Helsinki"},
			CorrectAnswers:   []string{"Stockholm"},
			Index:            0,
			TimeLimitSeconds: 20,
		},
		{
			ID:               "questionID2",
			QuizID:           quizID,
			Question:         "What is the capital of Norway?",
			Answers:          []string{"Stockholm", "Oslo", "Copenhagen", "Helsinki"},
			CorrectAnswers:   []string{"Oslo"},
			Index:            1,
			TimeLimitSeconds: 10,
		},
	}, nil)

	req, err := http.NewRequest("GET", server.URL+"/quizzes/"+quizID, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var quiz struct {
		Quiz      quizzer.Quiz       `json:"quiz"`
		Questions []quizzer.Question `json:"questions"`
	}
	err = json.NewDecoder(resp.Body).Decode(&quiz)
	require.NoError(t, err)
	require.Equal(t, "testquiz", quiz.Quiz.Title)
	require.Len(t, quiz.Questions, 2)
	require.Equal(t, "What is the capital of Sweden?", quiz.Questions[0].Question)
	require.Equal(t, "What is the capital of Norway?", quiz.Questions[1].Question)
}

func TestQuiz_ListQuizzes(t *testing.T) {
	client, sessionMock, server := authenticatedClient(t)
	defer server.Close()

	now := time.Now()
	testQuizzes := []quizzer.Quiz{
		{ID: "quizID1", Title: "testquiz1", CreatedBy: "testuser", CreatedAt: now},
		{ID: "quizID2", Title: "testquiz2", CreatedBy: "testuser", CreatedAt: now},
	}
	sessionMock.On("ListQuizzes", anyCtx).
		Return(testQuizzes, nil)

	req, err := http.NewRequest("GET", server.URL+"/quizzes", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var quizzes []quizzer.Quiz
	err = json.NewDecoder(resp.Body).Decode(&quizzes)
	require.NoError(t, err)

	for i := range quizzes {
		quizzes[i].CreatedAt = now
	}

	require.Equal(t, testQuizzes, quizzes)
}

func TestQuiz_DeleteQuiz(t *testing.T) {
	client, sessionMock, server := authenticatedClient(t)
	defer server.Close()

	t.Run("cannot delete another user's quiz", func(t *testing.T) {
		quizID := "quizID1"
		sessionMock.On("GetQuiz", anyCtx, quizID).
			Return(quizzer.Quiz{ID: quizID, CreatedBy: "test"}, nil)

		req, err := http.NewRequest("DELETE", server.URL+"/quizzes/"+quizID, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("delete own quiz", func(t *testing.T) {
		quizID := "quizID2"
		sessionMock.On("GetQuiz", anyCtx, quizID).
			Return(quizzer.Quiz{ID: quizID, CreatedBy: "test-userid"}, nil)

		sessionMock.On("DeleteQuiz", anyCtx, quizID).
			Return(nil)

		req, err := http.NewRequest("DELETE", server.URL+"/quizzes/"+quizID, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
