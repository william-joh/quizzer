package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	db := SetupTestDB(t)

	t.Run("get non-existing user", func(t *testing.T) {
		_, err := db.Do(context.Background()).GetUser(context.Background(), "testuser-id")
		require.Error(t, err)
	})

	t.Run("create user", func(t *testing.T) {
		err := db.Do(context.Background()).CreateUser(context.Background(), "testuser-id", "testuser", "testpassword")
		require.NoError(t, err)
	})

	t.Run("get existing user", func(t *testing.T) {
		user, err := db.Do(context.Background()).GetUser(context.Background(), "testuser-id")
		require.NoError(t, err)
		require.Equal(t, "testuser-id", user.ID)
		require.Equal(t, "testuser", user.Username)
		require.Equal(t, "testpassword", user.Password)
		require.NotZero(t, user.SignupDate)
	})

	t.Run("delete user", func(t *testing.T) {
		err := db.Do(context.Background()).DeleteUser(context.Background(), "testuser-id")
		require.NoError(t, err)
	})

	t.Run("get non-existing user", func(t *testing.T) {
		_, err := db.Do(context.Background()).GetUser(context.Background(), "testuser-id")
		require.Error(t, err)
	})
}
