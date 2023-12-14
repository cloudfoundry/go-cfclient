package config

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6InRlc3QgY2YgdG9rZW4iLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxNjIzOTAyMn0.mLvUvu-ED_lIkyI3UTXS_hUEPPFdI0BdNqRMgMThAhk"
const refreshToken = "secret-refresh-token"

func TestUsernamePassword(t *testing.T) {
	t.Run("with empty username", func(t *testing.T) {
		_, err := New("https://api.example.com",
			UserPassword("", "test"),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.Error(t, err)
		require.EqualError(t, err, "username and password are required when using using user credentials")
	})

	t.Run("with empty password", func(t *testing.T) {
		_, err := New("https://api.example.com", UserPassword("user", ""))
		require.Error(t, err)
		require.EqualError(t, err, "username and password are required when using using user credentials")
	})

	t.Run("with username and password", func(t *testing.T) {
		// user/pass hits the token endpoint
		uaaURL := testutil.SetupFakeUAAServer(300)
		c, err := New("https://api.example.com",
			UserPassword("username", "password"),
			AuthTokenURL(uaaURL, uaaURL))
		require.NoError(t, err)
		require.Equal(t, uaaURL, c.loginEndpointURL)
		require.Equal(t, uaaURL, c.uaaEndpointURL)
		require.Equal(t, "username", c.username)
		require.Equal(t, "password", c.password)
		require.Equal(t, "cf", c.clientID)
		require.Equal(t, GrantTypeAuthorizationCode, c.grantType)
	})

	t.Run("with username and password with non-default client", func(t *testing.T) {
		// user/pass hits the token endpoint
		uaaURL := testutil.SetupFakeUAAServer(300)
		c, err := New("https://api.example.com",
			UserPassword("username", "password"),
			ClientCredentials("clientID", ""),
			AuthTokenURL(uaaURL, uaaURL))
		require.NoError(t, err)
		require.Equal(t, "username", c.username)
		require.Equal(t, "password", c.password)
		require.Equal(t, "clientID", c.clientID)
		require.Equal(t, GrantTypeAuthorizationCode, c.grantType)
	})
}

func TestClientCredentials(t *testing.T) {
	t.Run("with invalid URL", func(t *testing.T) {
		_, err := New(":", ClientCredentials("clientID", "clientSecret"))
		require.ErrorContains(t, err, "expected an http(s) CF API root URI, but got")
	})

	t.Run("with clientID and empty client secret", func(t *testing.T) {
		_, err := New("https://api.example.com",
			ClientCredentials("clientID", ""),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.Error(t, err)
		require.EqualError(t, err, "CF API credentials were not provided")
	})

	t.Run("with empty clientID", func(t *testing.T) {
		c, err := New("https://api.example.com",
			ClientCredentials("", "clientSecret"),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.Equal(t, "cf", c.clientID)
		require.Equal(t, "clientSecret", c.clientSecret)
		require.Equal(t, GrantTypeClientCredentials, c.grantType)
	})

	t.Run("with clientID and client secret", func(t *testing.T) {
		c, err := New("https://api.example.com",
			ClientCredentials("clientID", "clientSecret"),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.Equal(t, "clientID", c.clientID)
		require.Equal(t, "clientSecret", c.clientSecret)
		require.Equal(t, GrantTypeClientCredentials, c.grantType)
	})

	t.Run("with clientID and client secret and access token", func(t *testing.T) {
		c, err := New("https://api.example.com",
			ClientCredentials("clientID", "clientSecret"),
			Token(accessToken, refreshToken),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.Equal(t, "clientID", c.clientID)
		require.Equal(t, "clientSecret", c.clientSecret)
		require.NotNil(t, c.oAuthToken)
		require.Equal(t, accessToken, c.oAuthToken.AccessToken)
		require.Equal(t, refreshToken, c.oAuthToken.RefreshToken)
		require.Equal(t, GrantTypeClientCredentials, c.grantType)
	})
}

func TestToken(t *testing.T) {
	t.Run("with empty token", func(t *testing.T) {
		_, err := New("https://api.example.com",
			Token("", ""),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid CF API token")
	})

	t.Run("with access token", func(t *testing.T) {
		c, err := New("https://api.example.com",
			Token(accessToken, ""),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.NotNil(t, c.oAuthToken)
		require.Equal(t, accessToken, c.oAuthToken.AccessToken)
		require.Equal(t, "", c.oAuthToken.RefreshToken)
		require.Equal(t, GrantTypeRefreshToken, c.grantType)
	})

	t.Run("with refresh token", func(t *testing.T) {
		c, err := New("https://api.example.com",
			Token("", refreshToken),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.NotNil(t, c.oAuthToken)
		require.Equal(t, GrantTypeRefreshToken, c.grantType)
	})

	t.Run("with access token and refresh token", func(t *testing.T) {
		c, err := New("https://api.example.com",
			Token(accessToken, refreshToken),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.NotNil(t, c.oAuthToken)
		require.Equal(t, GrantTypeRefreshToken, c.grantType)
	})

	t.Run("with access token and custom clientID", func(t *testing.T) {
		c, err := New("https://api.example.com",
			Token(accessToken, ""),
			ClientCredentials("myapp", ""),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com")) // skip service discovery
		require.NoError(t, err)
		require.Equal(t, GrantTypeRefreshToken, c.grantType)
	})

	t.Run("with custom http.Client without a transport and skip TLS verification", func(t *testing.T) {
		c, err := New("https://api.example.com",
			Token(accessToken, refreshToken),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com"), // skip service discovery
			SkipTLSValidation(),
			HttpClient(&http.Client{}))
		require.NoError(t, err)
		require.NotNil(t, c.HTTPClient())
		require.NotNil(t, c.HTTPClient().Transport)
		require.NotNil(t, c.HTTPClient().Transport.(*http.Transport).TLSClientConfig)
		require.True(t, c.HTTPClient().Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify)
	})
}

func TestNewConfigFromCFHomeDir(t *testing.T) {
	cfHomeDir := writeTestCFCLIConfig(t)

	t.Run("without overrides", func(t *testing.T) {
		cfg, err := NewFromCFHomeDir(cfHomeDir)
		require.NoError(t, err)
		require.Equal(t, "https://api.sys.example.com", cfg.apiEndpointURL)
		require.Equal(t, "https://login.sys.example.com", cfg.loginEndpointURL)
		require.Equal(t, "https://uaa.sys.example.com", cfg.uaaEndpointURL)
		require.Equal(t, accessToken, cfg.oAuthToken.AccessToken)
		require.Equal(t, refreshToken, cfg.oAuthToken.RefreshToken)
		require.Equal(t, DefaultClientID, cfg.clientID)
		require.Equal(t, GrantTypeRefreshToken, cfg.grantType)
	})

	t.Run("with with CF_USERNAME and CF_PASSWORD set", func(t *testing.T) {
		uaaURL := testutil.SetupFakeUAAServer(300)
		require.NoError(t, os.Setenv("CF_USERNAME", "admin"))
		require.NoError(t, os.Setenv("CF_PASSWORD", "pass"))
		defer func() {
			_ = os.Unsetenv("CF_USERNAME")
			_ = os.Unsetenv("CF_PASSWORD")
		}()
		cfg, err := NewFromCFHomeDir(cfHomeDir, AuthTokenURL(uaaURL, uaaURL))
		require.NoError(t, err)
		require.Equal(t, "https://api.sys.example.com", cfg.apiEndpointURL)
		require.Equal(t, uaaURL, cfg.loginEndpointURL)
		require.Equal(t, uaaURL, cfg.uaaEndpointURL)
		require.Equal(t, "admin", cfg.username)
		require.Equal(t, "pass", cfg.password)
		require.Equal(t, DefaultClientID, cfg.clientID)
		require.Equal(t, GrantTypeAuthorizationCode, cfg.grantType)
	})

	t.Run("with override options", func(t *testing.T) {
		uaaURL := testutil.SetupFakeUAAServer(300)
		cfg, err := NewFromCFHomeDir(cfHomeDir,
			UserPassword("admin", "pass"),
			AuthTokenURL(uaaURL, uaaURL))
		require.NoError(t, err)
		require.Equal(t, "https://api.sys.example.com", cfg.apiEndpointURL)
		require.Equal(t, uaaURL, cfg.loginEndpointURL)
		require.Equal(t, uaaURL, cfg.uaaEndpointURL)
		require.Equal(t, "admin", cfg.username)
		require.Equal(t, "pass", cfg.password)
		require.Equal(t, DefaultClientID, cfg.clientID)
		require.Equal(t, GrantTypeAuthorizationCode, cfg.grantType)
	})
}
