package api_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func TestQuiz_CreateQuiz(t *testing.T) {
	client, sessionMock, server := authenticatedClient(t)
	defer server.Close()

	var quizID string
	quiz := `{"title": "testquiz"}`
	sessionMock.On("CreateQuiz", anyCtx, anyString, "testquiz", anyString).
		Return(nil)

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
}

func TestQuiz_GetQuiz(t *testing.T) {
	client, sessionMock, server := authenticatedClient(t)
	defer server.Close()

	quizID := "quizID"
	sessionMock.On("GetQuiz", anyCtx, quizID).
		Return(quizzer.Quiz{ID: quizID, Title: "testquiz"}, nil)

	req, err := http.NewRequest("GET", server.URL+"/quizzes/"+quizID, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var quiz quizzer.Quiz
	err = json.NewDecoder(resp.Body).Decode(&quiz)
	require.NoError(t, err)
	require.Equal(t, "testquiz", quiz.Title)
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
