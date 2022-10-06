package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestListStacksByQuery(t *testing.T) {
	mocks := []MockRoute{
		{"GET", "/v3/stacks", []string{listStacksPayload}, "", http.StatusOK, "", nil},
	}
	setupMultiple(mocks, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	stacks, err := client.Stacks.ListByQuery(nil)
	require.NoError(t, err)
	require.Len(t, stacks, 2)

	require.Equal(t, "my-stack-1", stacks[0].Name)
	require.Equal(t, "This is my first stack!", stacks[0].Description)
	require.Equal(t, "guid-1", stacks[0].GUID)
	require.Equal(t, "https://api.example.org/v3/stacks/guid-1", stacks[0].Links["self"].Href)
	require.Equal(t, "my-stack-2", stacks[1].Name)
	require.Equal(t, "This is my second stack!", stacks[1].Description)
	require.Equal(t, "guid-2", stacks[1].GUID)
	require.Equal(t, "https://api.example.org/v3/stacks/guid-2", stacks[1].Links["self"].Href)
}
