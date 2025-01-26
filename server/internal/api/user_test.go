package api_test

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// 	"github.com/william-joh/quizzer/server/internal/quizzer"
// )

// var (
// 	anyCtx    = mock.MatchedBy(func(ctx context.Context) bool { return true })
// 	anyString = mock.MatchedBy(func(s string) bool { return true })
// )

// func TestUser(t *testing.T) {
// 	client, sessionMock, server := authenticatedClient(t)
// 	defer server.Close()

// 	user := `{"username": "testuser", "password": "testpassword"}`

// 	var userID string
// 	t.Run("create user", func(t *testing.T) {
// 		sessionMock.On("CreateUser", anyCtx, anyString, "testuser", anyString).Return(nil)
// 		sessionMock.On("CreateAuthSession", anyCtx, anyString, anyString).
// 			Return(nil)

// 		req, err := http.NewRequest("POST", server.URL+"/users", strings.NewReader(user))
// 		require.NoError(t, err)

// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		require.Equal(t, http.StatusCreated, resp.StatusCode)

// 		// Get userId from body
// 		var respBody struct {
// 			ID string `json:"id"`
// 		}
// 		err = json.NewDecoder(resp.Body).Decode(&respBody)
// 		require.NoError(t, err)
// 		userID = respBody.ID
// 		require.NotEmpty(t, userID)
// 	})

// 	t.Run("get user", func(t *testing.T) {
// 		exampleUser := quizzer.User{
// 			ID:         userID,
// 			Username:   "testuser",
// 			SignupDate: time.Now(),
// 		}
// 		sessionMock.On("GetUser", anyCtx, "test-userid").Return(exampleUser, nil)

// 		req, err := http.NewRequest("GET", server.URL+"/current-user", nil)
// 		require.NoError(t, err)

// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		require.Equal(t, http.StatusOK, resp.StatusCode)
// 		require.Equal(t, "application/json", resp.Header.Get("Content-Type"))

// 		var respUser quizzer.User
// 		err = json.NewDecoder(resp.Body).Decode(&respUser)
// 		require.NoError(t, err)
// 		respUser.SignupDate = exampleUser.SignupDate
// 		require.Equal(t, exampleUser, respUser)
// 	})

// 	t.Run("delete user", func(t *testing.T) {
// 		sessionMock.On("DeleteUser", anyCtx, userID).Return(nil)

// 		req, err := http.NewRequest("DELETE", server.URL+"/users/"+userID, nil)
// 		require.NoError(t, err)

// 		resp, err := client.Do(req)
// 		require.NoError(t, err)
// 		require.Equal(t, http.StatusOK, resp.StatusCode)
// 	})

// 	sessionMock.AssertExpectations(t)
// }
