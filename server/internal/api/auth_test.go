package api_test

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/william-joh/quizzer/server/internal/api"
	"github.com/william-joh/quizzer/server/internal/mocks"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func TestAuth_InvalidCredentials(t *testing.T) {
	client := http.Client{}
	dbMock := mocks.Database{}
	sessionMock := mocks.Session{}
	dbMock.On("Do", anyCtx).Return(&sessionMock)

	sessionMock.On("GetUserByUsername", anyCtx, "baduser").
		Return(quizzer.User{ID: "bad-userid", Username: "baduser", Password: "badpassword"}, nil)

	server := httptest.NewServer(api.NewAPI(&dbMock).Handler())
	defer server.Close()

	user := `{"username": "baduser", "password": "testpassword"}`
	req, err := http.NewRequest(http.MethodPost, server.URL+"/auth", strings.NewReader(user))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_UserNotFound(t *testing.T) {
	client := http.Client{}
	dbMock := mocks.Database{}
	sessionMock := mocks.Session{}
	dbMock.On("Do", anyCtx).Return(&sessionMock)

	sessionMock.On("GetUserByUsername", anyCtx, "baduser").
		Return(quizzer.User{}, errors.New("not found"))

	server := httptest.NewServer(api.NewAPI(&dbMock).Handler())
	defer server.Close()

	user := `{"username": "baduser", "password": "testpassword"}`
	req, err := http.NewRequest(http.MethodPost, server.URL+"/auth", strings.NewReader(user))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuth(t *testing.T) {
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := http.Client{
		Jar: jar,
	}
	user := `{"username": "testuser", "password": "testpassword"}`

	dbMock := mocks.Database{}
	sessionMock := mocks.Session{}
	dbMock.On("Do", anyCtx).Return(&sessionMock)

	sessionMock.On("GetUserByUsername", anyCtx, "testuser").
		Return(quizzer.User{ID: "test-userid", Username: "testuser", Password: "$2a$10$9uE.VDznFPnNzJx97ujxM.xoW6sUxcSe.s3bDazeRwHSL8nDzf2SK"}, nil)

	sessionMock.On("CreateAuthSession", anyCtx, "test-userid", anyString).
		Return(nil)

	sessionMock.On("GetAuthSession", anyCtx, anyString).
		Return("test-userid", nil)

	sessionMock.On("GetUser", anyCtx, "test-userid").
		Return(quizzer.User{ID: "test-userid", Username: "testuser"}, nil)

	server := httptest.NewServer(api.NewAPI(&dbMock).Handler())
	defer server.Close()

	t.Run("unauthenticated client cannot access user routes", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/users/test-userid", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("authenticate user", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, server.URL+"/auth", strings.NewReader(user))
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		// check if cookie is set
		require.NotEmpty(t, resp.Cookies())
		require.Equal(t, "quizzer_session_id", resp.Cookies()[0].Name)

		// subsequent requests should be authenticated by the cookie
		req, err = http.NewRequest(http.MethodGet, server.URL+"/users/test-userid", nil)
		require.NoError(t, err)

		resp, err = client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func authenticatedClient(t *testing.T) (http.Client, *mocks.Session, *httptest.Server) {
	t.Helper()
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := http.Client{
		Jar: jar,
	}

	user := `{"username": "testuser", "password": "testpassword"}`

	dbMock := mocks.Database{}
	sessionMock := mocks.Session{}
	dbMock.On("Do", anyCtx).Return(&sessionMock)

	sessionMock.On("GetUserByUsername", anyCtx, "testuser").
		Return(quizzer.User{ID: "test-userid", Username: "testuser", Password: "$2a$10$9uE.VDznFPnNzJx97ujxM.xoW6sUxcSe.s3bDazeRwHSL8nDzf2SK"}, nil)

	sessionMock.On("CreateAuthSession", anyCtx, "test-userid", anyString).
		Return(nil)

	sessionMock.On("GetAuthSession", anyCtx, anyString).
		Return("test-userid", nil).Maybe()

	server := httptest.NewServer(api.NewAPI(&dbMock).Handler())

	req, err := http.NewRequest(http.MethodPost, server.URL+"/auth", strings.NewReader(user))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	// check if cookie is set
	require.NotEmpty(t, resp.Cookies())
	require.Equal(t, "quizzer_session_id", resp.Cookies()[0].Name)

	return client, &sessionMock, server
}
