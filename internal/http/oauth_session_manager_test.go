package http_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"testing"
)

func TestOAuthSessionManager(t *testing.T) {
	uaaURL := testutil.SetupFakeUAAServer(300)
	defer testutil.Teardown()

	// missing configuration
	c, err := config.NewUserPassword("https://api.example.org", "admin", "secret")
	require.NoError(t, err)
	require.Empty(t, c.LoginEndpointURL)
	require.Empty(t, c.UAAEndpointURL)
	m := http.NewOAuthSessionManager(c)

	_, err = m.Client(true)
	require.Error(t, err, "expected an error when UAA or Login endpoint is empty")
	require.Equal(t, "login and UAA endpoints must not be empty", err.Error())

	// minimal proper configuration
	c.LoginEndpointURL = uaaURL
	c.UAAEndpointURL = uaaURL

	// we can create a client that utilizes oauth
	client1, err := m.Client(true)
	require.NoError(t, err)
	require.NotNil(t, client1)

	// the same access token is returned as long as it's not expired (which it's not - 300s)
	token, err := m.AccessToken()
	require.NoError(t, err)
	require.Equal(t, "foobar1", token)
	require.NoError(t, err)
	token, err = m.AccessToken()
	require.NoError(t, err)
	require.Equal(t, "foobar1", token)

	// the same client is returned
	client2, err := m.Client(true)
	require.NoError(t, err)
	require.Same(t, client1, client2)

	// we force new auth context
	err = m.ReAuthenticate()
	require.NoError(t, err)

	// a different client is now returned
	client3, err := m.Client(true)
	require.NoError(t, err)
	require.NotSame(t, client2, client3)

	// a new token is also returned
	token, err = m.AccessToken()
	require.NoError(t, err)
	require.Equal(t, "foobar2", token)

	// however, it still reuses the client's underlying transport as that's what pools TCP connections
	// even after re-establishing auth context
	oauthTransport1 := client2.Transport.(*oauth2.Transport)
	oauthTransport2 := client3.Transport.(*oauth2.Transport)

	require.Same(t, oauthTransport1.Base, oauthTransport2.Base, "expect the same http transport between Client() calls")
	require.Same(t, c.HTTPClient().Transport, oauthTransport2.Base, "expect the same http transport from config")
}
