package http_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOAuthSessionManager(t *testing.T) {
	uaaURL := testutil.SetupFakeUAAServer(3)
	defer testutil.Teardown()

	// missing configuration
	c, err := config.NewUserPassword("https://api.example.org", "admin", "secret")
	require.NoError(t, err)
	require.Empty(t, c.LoginEndpointURL)
	require.Empty(t, c.UAAEndpointURL)
	m := http.NewOAuthSessionManager(c)

	_, err = m.Client()
	require.Error(t, err, "expected an error when UAA or Login endpoint is empty")
	require.Equal(t, "login and UAA endpoints must not be empty", err.Error())

	// minimal proper configuration
	c.LoginEndpointURL = uaaURL
	c.UAAEndpointURL = uaaURL

	client, err := m.Client()
	require.NoError(t, err)
	require.NotNil(t, client)

	token, err := m.Token()
	require.NoError(t, err)
	require.Equal(t, "bearer foobar2", token)
	require.NoError(t, err)
	token, err = m.Token()
	require.Equal(t, "bearer foobar3", token)

	// it reuses the client transport as that's what re-uses TCP connections
	client2, err := m.Client()
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Same(t, client.Transport, client2.Transport, "expect the same http transport")

	// it can refresh when expired
	//time.Sleep(3 * time.Second)
	//client3, err := m.Client()
	//require.NoError(t, err)
	//require.NotNil(t, client)
	//require.NotSame(t, client, client3, "expect a new http client when expired")
}
