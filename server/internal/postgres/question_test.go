package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func TestQuestions(t *testing.T) {
	db := setupTestDB(t)

	err := db.Do(context.Background()).CreateUser(context.Background(), "testuser-id", "testuser", "testpassword")
	require.NoError(t, err)
	user, err := db.Do(context.Background()).GetUser(context.Background(), "testuser")
	require.NoError(t, err)

	err = db.Do(context.Background()).CreateQuiz(context.Background(), "testquiz-id", "testquiz", user.ID)
	require.NoError(t, err)

	t.Run("get non-existing question", func(t *testing.T) {
		_, err := db.Do(context.Background()).GetQuestion(context.Background(), "testquestion")
		require.Error(t, err)
	})

	t.Run("create question", func(t *testing.T) {
		err := db.Do(context.Background()).CreateQuestion(context.Background(), quizzer.Question{
			ID:             "testquestion-id1",
			QuizID:         "testquiz-id",
			Question:       "testquestion1",
			Index:          1,
			TimeLimit:      10 * time.Second,
			Answers:        []string{"answer1", "answer2", "answer3"},
			CorrectAnswer:  "answer1",
			VideoURL:       asPtr("testurl"),
			VideoStartTime: asPtr(time.Millisecond * 100),
			VideoEndTime:   asPtr(time.Millisecond * 200),
		})
		require.NoError(t, err)

		err = db.Do(context.Background()).CreateQuestion(context.Background(), quizzer.Question{
			ID:            "testquestion-id2",
			QuizID:        "testquiz-id",
			Question:      "testquestion2",
			Index:         2,
			TimeLimit:     20 * time.Second,
			Answers:       []string{"answer1", "answer2", "answer3"},
			CorrectAnswer: "answer2",
		})
		require.NoError(t, err)

		err = db.Do(context.Background()).CreateQuestion(context.Background(), quizzer.Question{
			ID:            "testquestion-id3",
			QuizID:        "testquiz-id",
			Question:      "testquestion3",
			Index:         3,
			TimeLimit:     30 * time.Second,
			Answers:       []string{"answer1", "answer2", "answer3"},
			CorrectAnswer: "answer3",
		})
		require.NoError(t, err)
	})

	t.Run("list questions", func(t *testing.T) {
		questions, err := db.Do(context.Background()).ListQuestions(context.Background(), "testquiz-id")
		require.NoError(t, err)
		require.Len(t, questions, 3)

		expectedQuestion1 := quizzer.Question{
			ID:             "testquestion-id1",
			QuizID:         "testquiz-id",
			Question:       "testquestion1",
			Index:          1,
			TimeLimit:      10 * time.Second,
			Answers:        []string{"answer1", "answer2", "answer3"},
			CorrectAnswer:  "answer1",
			VideoURL:       asPtr("testurl"),
			VideoStartTime: asPtr(time.Millisecond * 100),
			VideoEndTime:   asPtr(time.Millisecond * 200),
		}
		require.Equal(t, expectedQuestion1, questions[0])

		expectedQuestion2 := quizzer.Question{
			ID:            "testquestion-id2",
			QuizID:        "testquiz-id",
			Question:      "testquestion2",
			Index:         2,
			TimeLimit:     20 * time.Second,
			Answers:       []string{"answer1", "answer2", "answer3"},
			CorrectAnswer: "answer2",
			VideoURL:      nil,
		}
		require.Equal(t, expectedQuestion2, questions[1])
	})

	t.Run("get question", func(t *testing.T) {
		question, err := db.Do(context.Background()).GetQuestion(context.Background(), "testquestion-id1")
		require.NoError(t, err)
		expectedQuestion := quizzer.Question{
			ID:             "testquestion-id1",
			QuizID:         "testquiz-id",
			Question:       "testquestion1",
			Index:          1,
			TimeLimit:      10 * time.Second,
			Answers:        []string{"answer1", "answer2", "answer3"},
			CorrectAnswer:  "answer1",
			VideoURL:       asPtr("testurl"),
			VideoStartTime: asPtr(time.Millisecond * 100),
			VideoEndTime:   asPtr(time.Millisecond * 200),
		}
		require.Equal(t, expectedQuestion, question)
	})

	t.Run("edit question", func(t *testing.T) {
		err := db.Do(context.Background()).UpdateQuestion(context.Background(), quizzer.Question{
			ID:             "testquestion-id1",
			QuizID:         "testquiz-id",
			Question:       "editedquestion1",
			Index:          1,
			TimeLimit:      10 * time.Second,
			Answers:        []string{"answer1", "answer2", "answer3", "answer4"},
			CorrectAnswer:  "answer3",
			VideoURL:       asPtr("editedurl"),
			VideoStartTime: asPtr(time.Millisecond * 200),
			VideoEndTime:   asPtr(time.Millisecond * 300),
		})
		require.NoError(t, err)

		question, err := db.Do(context.Background()).GetQuestion(context.Background(), "testquestion-id1")
		require.NoError(t, err)
		expectedQuestion := quizzer.Question{
			ID:             "testquestion-id1",
			QuizID:         "testquiz-id",
			Question:       "editedquestion1",
			Index:          1,
			TimeLimit:      10 * time.Second,
			Answers:        []string{"answer1", "answer2", "answer3", "answer4"},
			CorrectAnswer:  "answer3",
			VideoURL:       asPtr("editedurl"),
			VideoStartTime: asPtr(time.Millisecond * 200),
			VideoEndTime:   asPtr(time.Millisecond * 300),
		}
		require.Equal(t, expectedQuestion, question)
	})

	t.Run("delete question", func(t *testing.T) {
		err := db.Do(context.Background()).DeleteQuestion(context.Background(), "testquestion-id2")
		require.NoError(t, err)

		questions, err := db.Do(context.Background()).ListQuestions(context.Background(), "testquiz-id")
		require.NoError(t, err)
		require.Len(t, questions, 2)

		require.Equal(t, "testquestion-id1", questions[0].ID)
		require.Equal(t, "testquestion-id3", questions[1].ID)
	})
}

func asPtr[T any](s T) *T {
	return &s
}
