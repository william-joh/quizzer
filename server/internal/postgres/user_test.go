package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)

	t.Run("get non-existing user", func(t *testing.T) {
		_, err := db.Do(context.Background()).GetUser(context.Background(), "testuser")
		require.Error(t, err)
	})

	t.Run("create user", func(t *testing.T) {
		err := db.Do(context.Background()).CreateUser(context.Background(), "testuser", "testpassword")
		require.NoError(t, err)
	})

	t.Run("get existing user", func(t *testing.T) {
		user, err := db.Do(context.Background()).GetUser(context.Background(), "testuser")
		require.NoError(t, err)
		require.Equal(t, "testuser", user.Username)
		require.Equal(t, "testpassword", user.Password)
		require.NotEmpty(t, user.ID)
		require.NotZero(t, user.SignupDate)
	})

	t.Run("delete user", func(t *testing.T) {
		err := db.Do(context.Background()).DeleteUser(context.Background(), "testuser")
		require.NoError(t, err)
	})

	t.Run("get non-existing user", func(t *testing.T) {
		_, err := db.Do(context.Background()).GetUser(context.Background(), "testuser")
		require.Error(t, err)
	})
}
