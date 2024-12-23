package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	session := setupTestDB(t)

	t.Run("get non-existing user", func(t *testing.T) {
		_, err := session.GetUser(context.Background(), "testuser")
		require.Error(t, err)
	})

	t.Run("create user", func(t *testing.T) {
		err := session.CreateUser(context.Background(), "testuser", "testpassword")
		fmt.Println("user created")
		require.NoError(t, err)
	})

	t.Run("get existing user", func(t *testing.T) {
		user, err := session.GetUser(context.Background(), "testuser")
		require.NoError(t, err)
		require.Equal(t, "testuser", user.Username)
		require.Equal(t, "testpassword", user.Password)
		require.NotEmpty(t, user.ID)
		require.NotZero(t, user.SignupDate)
	})

	t.Run("delete user", func(t *testing.T) {
		err := session.DeleteUser(context.Background(), "testuser")
		require.NoError(t, err)
	})

	t.Run("get non-existing user", func(t *testing.T) {
		_, err := session.GetUser(context.Background(), "testuser")
		require.Error(t, err)
	})
}
