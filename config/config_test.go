package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		c, err := New("https://api.example.com", Token("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6InRlc3QgY2YgdG9rZW4iLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxNjIzOTAyMn0.mLvUvu-ED_lIkyI3UTXS_hUEPPFdI0BdNqRMgMThAhk", "test"),
			AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com"))
		require.Nil(t, err)
		require.NotNil(t, c.oAuthToken)
		require.Equal(t, GrantTypeRefreshToken, c.grantType)
	})

	t.Run("Test with valid ClientCredentials", func(t *testing.T) {
		c, err := New("https://api.example.com", ClientCredentials("clientID", "clientSecret"), AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com"))
		require.Nil(t, err)
		require.Equal(t, "https://login.cf.example.com", c.loginEndpointURL)
		require.Equal(t, "https://token.cf.example.com", c.uaaEndpointURL)
		require.Equal(t, GrantTypeClientCredentials, c.grantType)
	})

	t.Run("Test with valid UserPassword", func(t *testing.T) {
		c, err := New("https://api.example.com", UserPassword("username", "password"), AuthTokenURL("https://login.cf.example.com", "https://token.cf.example.com"))
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
	cfHomeDir, err := os.MkdirTemp("", "cf_home")
	require.NoError(t, err)

	configDir := path.Join(cfHomeDir, ".cf")
	err = os.MkdirAll(configDir, 0744)
	require.NoError(t, err)

	configPath := path.Join(configDir, "config.json")
	err = os.WriteFile(configPath, []byte(cfCLIConfig), 0744)
	require.NoError(t, err)

	cfg, err := NewFromCFHomeDir(cfHomeDir)
	require.NoError(t, err)
	require.Equal(t, "https://api.sys.example.com", cfg.apiEndpointURL)
	require.Equal(t, DefaultClientID, cfg.clientID)
	require.Equal(t, "https://uaa.sys.example.com", cfg.uaaEndpointURL)
	require.Equal(t, GrantTypeRefreshToken, cfg.grantType)
}

const cfCLIConfig = `
{
  "ConfigVersion": 3,
  "Target": "https://api.sys.example.com",
  "APIVersion": "2.164.0",
  "AuthorizationEndpoint": "https://login.sys.example.com",
  "DopplerEndPoint": "wss://doppler.sys.example.com:443",
  "UaaEndpoint": "https://uaa.sys.example.com",
  "RoutingAPIEndpoint": "https://api.sys.example.com/routing",
  "AccessToken": "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6InRlc3QgY2YgdG9rZW4iLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxNjIzOTAyMn0.mLvUvu-ED_lIkyI3UTXS_hUEPPFdI0BdNqRMgMThAhk",
  "SSHOAuthClient": "ssh-proxy",
  "UAAOAuthClient": "cf",
  "UAAOAuthClientSecret": "",
  "UAAGrantType": "",
  "RefreshToken": "secret-refresh-token",
  "OrganizationFields": {
    "GUID": "42754be1-f558-4d28-9c06-c706f6641245",
    "Name": "system",
    "QuotaDefinition": {
      "guid": "",
      "name": "",
      "memory_limit": 0,
      "instance_memory_limit": 0,
      "total_routes": 0,
      "total_services": 0,
      "non_basic_services_allowed": false,
      "app_instance_limit": 0,
      "total_reserved_route_ports": 0
    }
  },
  "SpaceFields": {
    "GUID": "e42ccfe9-04bf-4cbc-ae16-f26741778a71",
    "Name": "system",
    "AllowSSH": true
  },
  "SSLDisabled": true,
  "AsyncTimeout": 0,
  "Trace": "",
  "ColorEnabled": "",
  "Locale": "",
  "PluginRepos": [
    {
      "Name": "CF-Community",
      "URL": "https://plugins.cloudfoundry.org"
    }
  ],
  "MinCLIVersion": "6.23.0",
  "MinRecommendedCLIVersion": "6.23.0"
}
`
