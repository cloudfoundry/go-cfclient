package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadCFCLIConfig(t *testing.T) {
	cfHomeDir := writeTestCFCLIConfig(t)
	cf, err := loadCFCLIConfig(cfHomeDir)
	require.NoError(t, err)
	require.Equal(t, "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6InRlc3QgY2YgdG9rZW4iLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxNjIzOTAyMn0.mLvUvu-ED_lIkyI3UTXS_hUEPPFdI0BdNqRMgMThAhk", cf.AccessToken)
	require.Equal(t, "secret-refresh-token", cf.RefreshToken)
	require.Equal(t, "https://api.sys.example.com", cf.Target)
	require.Equal(t, "https://login.sys.example.com", cf.AuthorizationEndpoint)
	require.Equal(t, "https://uaa.sys.example.com", cf.UaaEndpoint)
	require.Equal(t, "cf", cf.UAAOAuthClient)
	require.Equal(t, "", cf.UAAOAuthClientSecret)
	require.Equal(t, "", cf.UAAGrantType)
	require.Equal(t, "ssh-proxy", cf.SSHOAuthClient)
	require.True(t, cf.SSLDisabled)
}

func writeTestCFCLIConfig(t *testing.T) string {
	cfHomeDir, err := os.MkdirTemp("", "cf_home")
	require.NoError(t, err)

	configDir := path.Join(cfHomeDir, ".cf")
	err = os.MkdirAll(configDir, 0744)
	require.NoError(t, err)

	configPath := path.Join(configDir, "config.json")
	err = os.WriteFile(configPath, []byte(cfCLIConfigJSON), 0744)
	require.NoError(t, err)

	return cfHomeDir
}

const cfCLIConfigJSON = `
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
