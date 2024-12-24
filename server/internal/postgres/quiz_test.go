package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func TestQuiz(t *testing.T) {
	db := setupTestDB(t)

	err := db.Do(context.Background()).CreateUser(context.Background(), "testuser-id", "testuser", "testpassword")
	require.NoError(t, err)
	user, err := db.Do(context.Background()).GetUser(context.Background(), "testuser")
	require.NoError(t, err)

	t.Run("get non-existing quiz", func(t *testing.T) {
		_, err := db.Do(context.Background()).GetQuiz(context.Background(), "testquiz")
		require.Error(t, err)
	})

	t.Run("create quiz", func(t *testing.T) {
		err := db.Do(context.Background()).CreateQuiz(context.Background(), "testquiz-id1", "testquiz1", user.ID)
		require.NoError(t, err)

		err = db.Do(context.Background()).CreateQuiz(context.Background(), "testquiz-id2", "testquiz2", user.ID)
		require.NoError(t, err)

		err = db.Do(context.Background()).CreateQuiz(context.Background(), "testquiz-id3", "testquiz3", user.ID)
		require.NoError(t, err)
	})

	var quizzes []quizzer.Quiz
	t.Run("list quizzes", func(t *testing.T) {
		quizzes, err = db.Do(context.Background()).ListQuizzes(context.Background())
		require.NoError(t, err)
		require.Len(t, quizzes, 3)

		require.Equal(t, "testquiz-id1", quizzes[0].ID)
		require.Equal(t, "testquiz1", quizzes[0].Title)
		require.Equal(t, user.ID, quizzes[0].CreatedBy)
		require.NotZero(t, quizzes[0].CreatedAt)

		require.Equal(t, "testquiz-id2", quizzes[1].ID)
		require.Equal(t, "testquiz2", quizzes[1].Title)
		require.Equal(t, user.ID, quizzes[1].CreatedBy)
		require.NotZero(t, quizzes[1].CreatedAt)
	})

	t.Run("get quiz", func(t *testing.T) {
		quiz, err := db.Do(context.Background()).GetQuiz(context.Background(), quizzes[0].ID)
		require.NoError(t, err)
		require.Equal(t, quizzes[0], quiz)
	})

	t.Run("delete quiz", func(t *testing.T) {
		err := db.Do(context.Background()).DeleteQuiz(context.Background(), quizzes[1].ID)
		require.NoError(t, err)

		quizzes, err = db.Do(context.Background()).ListQuizzes(context.Background())
		require.NoError(t, err)
		require.Len(t, quizzes, 2)

		require.Equal(t, "testquiz1", quizzes[0].Title)
		require.Equal(t, "testquiz3", quizzes[1].Title)
	})
}
