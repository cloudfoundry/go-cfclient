package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6InRlc3QgY2YgdG9rZW4iLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxNjIzOTAyMn0.mLvUvu-ED_lIkyI3UTXS_hUEPPFdI0BdNqRMgMThAhk"
const refreshToken = "secret-refresh-token"

func TestConfig(t *testing.T) {
	t.Run("Test with empty ClientID", func(t *testing.T) {
		_, err := New("https://api.example.com", ClientCredentials("", "test"))
		require.NotNil(t, err)
		require.Equal(t, "expected a non-empty CF API clientID", err.Error())
	})

	t.Run("Test with empty Username", func(t *testing.T) {
		_, err := New("https://api.example.com", UserPassword("", "test"))
		require.NotNil(t, err)
		require.Equal(t, "expected a non-empty CF API username", err.Error())
	})

	t.Run("Test with empty Password", func(t *testing.T) {
		_, err := New("https://api.example.com", UserPassword("user", ""))
		require.NotNil(t, err)
		require.Equal(t, "expected a non-empty CF API password", err.Error())
	})

	t.Run("Test with empty Refresh Token", func(t *testing.T) {
		_, err := New("https://api.example.com", Token("test", ""))
		require.NotNil(t, err)
		require.ErrorContains(t, err, "invalid CF API token:")
	})

	t.Run("Test with Valid Tokens", func(t *testing.T) {
		c, err := New("https://api.example.com", Token(accessToken, refreshToken),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com"))
		require.Nil(t, err)
		require.NotNil(t, c.oAuthToken)
		require.Equal(t, GrantTypeRefreshToken, c.grantType)
	})

	t.Run("Test with valid ClientCredentials", func(t *testing.T) {
		c, err := New("https://api.example.com",
			ClientCredentials("clientID", "clientSecret"),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com"))
		require.Nil(t, err)
		require.Equal(t, "https://login.cf.example.com", c.loginEndpointURL)
		require.Equal(t, "https://token.cf.example.com", c.uaaEndpointURL)
		require.Equal(t, GrantTypeClientCredentials, c.grantType)
	})

	t.Run("Test with valid UserPassword", func(t *testing.T) {
		c, err := New("https://api.example.com",
			UserPassword("username", "password"),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com"))
		require.Nil(t, err)
		require.Equal(t, "https://login.cf.example.com", c.loginEndpointURL)
		require.Equal(t, "https://token.cf.example.com", c.uaaEndpointURL)
		require.Equal(t, GrantTypeAuthorizationCode, c.grantType)
	})

	t.Run("Test with invalid URL", func(t *testing.T) {
		_, err := New(":", ClientCredentials("clientID", "clientSecret"))
		require.ErrorContains(t, err, "expected an http(s) CF API root URI, but got")
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
		require.NoError(t, os.Setenv("CF_USERNAME", "admin"))
		require.NoError(t, os.Setenv("CF_PASSWORD", "pass"))
		defer func() {
			_ = os.Unsetenv("CF_USERNAME")
			_ = os.Unsetenv("CF_PASSWORD")
		}()
		cfg, err := NewFromCFHomeDir(cfHomeDir)
		require.NoError(t, err)
		require.Equal(t, "https://api.sys.example.com", cfg.apiEndpointURL)
		require.Equal(t, "https://login.sys.example.com", cfg.loginEndpointURL)
		require.Equal(t, "https://uaa.sys.example.com", cfg.uaaEndpointURL)
		require.Equal(t, "admin", cfg.username)
		require.Equal(t, "pass", cfg.password)
		require.Equal(t, DefaultClientID, cfg.clientID)
		require.Equal(t, GrantTypeAuthorizationCode, cfg.grantType)
	})

	t.Run("with override options", func(t *testing.T) {
		cfg, err := NewFromCFHomeDir(cfHomeDir,
			UserPassword("admin", "pass"),
			AuthTokenURL("https://login2.sys.example.com", "https://uaa2.sys.example.com"))
		require.NoError(t, err)
		require.Equal(t, "https://api.sys.example.com", cfg.apiEndpointURL)
		require.Equal(t, "https://login2.sys.example.com", cfg.loginEndpointURL)
		require.Equal(t, "https://uaa2.sys.example.com", cfg.uaaEndpointURL)
		require.Equal(t, "admin", cfg.username)
		require.Equal(t, "pass", cfg.password)
		require.Equal(t, DefaultClientID, cfg.clientID)
		require.Equal(t, GrantTypeAuthorizationCode, cfg.grantType)
	})
}
