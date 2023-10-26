package http_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"testing"
)

const accessToken = `bearer ignored.eyJqdGkiOiJhOGE5YTJjNDY5MzY0YzU3YmI2M2QxMWFiYzdhNjgzOSIsInN1YiI6IjJiNmMzM2ZlLTExZTItNGQwMi05OTNhLTdiNjQ5ZjhhMmI5YyIsInNjb3BlIjpbIm9wZW5pZCIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy53cml0ZSIsIm5ldHdvcmsud3JpdGUiLCJzY2ltLnJlYWQiLCJjbG91ZF9jb250cm9sbGVyLmFkbWluIiwidWFhLnVzZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm5ldHdvcmsuYWRtaW4iLCJkb3BwbGVyLmZpcmVob3NlIiwic2NpbS53cml0ZSJdLCJjbGllbnRfaWQiOiJjZiIsImNpZCI6ImNmIiwiYXpwIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiMmI2YzMzZmUtMTFlMi00ZDAyLTk5M2EtN2I2NDlmOGEyYjljIiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoiYWRtaW4iLCJlbWFpbCI6ImFkbWluIiwiYXV0aF90aW1lIjoxNjk4MDk2Mzc2LCJyZXZfc2lnIjoiZmNlMmY2MDAiLCJpYXQiOjE2OTgwOTY0MDgsImV4cCI6MTY5ODA5NjQ2OCwiaXNzIjoiaHR0cHM6Ly91YWEuc3lzLmgyby0yLTE5MTQ5Lmgyby52bXdhcmUuY29tL29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiYXVkIjpbImRvcHBsZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMiLCJvcGVuaWQiLCJjbG91ZF9jb250cm9sbGVyIiwicGFzc3dvcmQiLCJzY2ltIiwidWFhIiwibmV0d29yayIsImNmIl19.ignored`

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

func TestOAuthSessionManagerRefreshToken(t *testing.T) {
	uaaURL := testutil.SetupFakeUAAServer(30)
	defer testutil.Teardown()

	// create a config using an expired access token
	c, err := config.NewToken("https://api.example.org", accessToken, "refresh-token")
	require.NoError(t, err)
	c.LoginEndpointURL = uaaURL
	c.UAAEndpointURL = uaaURL
	m := http.NewOAuthSessionManager(c)

	// we can create a client that utilizes oauth
	client1, err := m.Client(true)
	require.NoError(t, err)
	require.NotNil(t, client1)

	// get the access token, it should have been auto-refreshed because the one we gave in the config was expired
	token, err := m.AccessToken()
	require.NoError(t, err)
	require.NotEqual(t, accessToken, token)
	require.Equal(t, "foobar1", token)

	// get the access token, should be the same as before
	token, err = m.AccessToken()
	require.NoError(t, err)
	require.Equal(t, "foobar1", token)

	// we cannot re-auth with only a refresh token (no credentials)
	err = m.ReAuthenticate()
	require.EqualError(t, err, "cannot reauthenticate user token auth type, check your access and/or refresh token expiration date")
}
