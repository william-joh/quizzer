package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	db := SetupTestDB(t)

	err := db.Do(context.Background()).CreateUser(context.Background(), "testuser-id", "testuser", "testpassword")
	require.NoError(t, err)

	t.Run("cannot create session with non-existing user", func(t *testing.T) {
		err := db.Do(context.Background()).CreateAuthSession(context.Background(), "non-existing-user", "testsession")
		require.Error(t, err)
	})

	t.Run("create session", func(t *testing.T) {
		err := db.Do(context.Background()).CreateAuthSession(context.Background(), "testuser-id", "testsession")
		require.NoError(t, err)
	})

	t.Run("get session", func(t *testing.T) {
		userID, err := db.Do(context.Background()).GetAuthSession(context.Background(), "testsession")
		require.NoError(t, err)
		require.Equal(t, "testuser-id", userID)
	})
}
