package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestListUserByQuery(t *testing.T) {
	mocks := []MockRoute{
		{"GET", "/v3/users", []string{listUsersPayload}, "", http.StatusOK, "", nil},
		{"GET", "/v3/userspage2", []string{listUsersPayloadPage2}, "", http.StatusOK, "page=2&per_page=2", nil},
	}
	setupMultiple(mocks, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	query := url.Values{}
	users, err := client.Users.ListByQuery(query)
	require.NoError(t, err)
	require.Len(t, users, 3)

	require.Equal(t, "smoke_tests", users[0].Username)
	require.Equal(t, "test1", users[1].Username)
	require.Equal(t, "test2", users[2].Username)
}
