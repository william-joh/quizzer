package api_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func TestQuestion_CreateQuestion(t *testing.T) {
	tests := []struct {
		name             string
		givenQuestion    string
		expectedQuestion quizzer.Question
	}{
		{
			name: "valid minimal question",
			givenQuestion: `{
				"quizId": "quizID",
				"question": "testquestion",
				"answers": ["a", "b"],
				"correctAnswers": ["a"],
				"index": 0,
				"timeLimitSeconds": 30
			}`,
			expectedQuestion: quizzer.Question{
				QuizID:           "quizID",
				Question:         "testquestion",
				Answers:          []string{"a", "b"},
				CorrectAnswers:   []string{"a"},
				Index:            0,
				TimeLimitSeconds: 30,
			},
		},
		{
			name: "valid question with video",
			givenQuestion: `{
				"quizId": "quizID",
				"question": "testquestion",
				"answers": ["a", "b", "c"],
				"correctAnswers": ["c"],
				"index": 1,
				"timeLimitSeconds": 60,
				"videoUrl": "https://example.com",
				"videoStartTimeSeconds": 10,
				"videoEndTimeSeconds": 20
			}`,
			expectedQuestion: quizzer.Question{
				QuizID:                "quizID",
				Question:              "testquestion",
				Answers:               []string{"a", "b", "c"},
				CorrectAnswers:        []string{"c"},
				Index:                 1,
				TimeLimitSeconds:      60,
				VideoURL:              asPtr("https://example.com"),
				VideoStartTimeSeconds: asPtr(uint64(10)),
				VideoEndTimeSeconds:   asPtr(uint64(20)),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, sessionMock, server := authenticatedClient(t)
			defer server.Close()

			sessionMock.On("GetQuiz", anyCtx, "quizID").
				Return(quizzer.Quiz{CreatedBy: "test-userid"}, nil)

			e := mock.MatchedBy(func(gotQuestion quizzer.Question) bool {
				test.expectedQuestion.ID = gotQuestion.ID
				return reflect.DeepEqual(test.expectedQuestion, gotQuestion)
			})

			sessionMock.On("CreateQuestion", anyCtx, e).
				Return(nil)

			req, err := http.NewRequest("POST", server.URL+"/quizzes/quizID/questions", strings.NewReader(test.givenQuestion))
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, resp.StatusCode)

			// Get userId from body
			var respBody struct {
				ID string `json:"id"`
			}
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			require.NoError(t, err)
			require.NotEmpty(t, respBody.ID)
		})
	}
}

func TestQuestion_ListQuestions(t *testing.T) {
	client, sessionMock, server := authenticatedClient(t)
	defer server.Close()

	questions := []quizzer.Question{
		{
			ID:               "questionID1",
			QuizID:           "quizID",
			Question:         "testquestion1",
			Index:            0,
			TimeLimitSeconds: 10,
			Answers:          []string{"answer1", "answer2", "answer3"},
			CorrectAnswers:   []string{"answer1"},
		},
		{
			ID:                    "questionID2",
			QuizID:                "quizID",
			Question:              "testquestion2",
			Index:                 1,
			TimeLimitSeconds:      20,
			Answers:               []string{"answer1", "answer2", "answer3"},
			CorrectAnswers:        []string{"answer2"},
			VideoURL:              asPtr("testurl"),
			VideoStartTimeSeconds: asPtr(uint64(10)),
			VideoEndTimeSeconds:   asPtr(uint64(20)),
		},
	}
	sessionMock.On("ListQuestions", anyCtx, "quizID").
		Return(questions, nil)

	req, err := http.NewRequest("GET", server.URL+"/quizzes/quizID/questions", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var gotQuestions []quizzer.Question
	err = json.NewDecoder(resp.Body).Decode(&gotQuestions)
	require.NoError(t, err)
	require.Equal(t, questions, gotQuestions)
}

func TestQuestion_DeleteQuestion(t *testing.T) {
	client, sessionMock, server := authenticatedClient(t)
	defer server.Close()

	sessionMock.On("GetQuiz", anyCtx, "quizID").
		Return(quizzer.Quiz{CreatedBy: "test-userid"}, nil)

	sessionMock.On("DeleteQuestion", anyCtx, "questionID").
		Return(nil)

	req, err := http.NewRequest("DELETE", server.URL+"/quizzes/quizID/questions/questionID", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func asPtr[T any](v T) *T {
	return &v
}
