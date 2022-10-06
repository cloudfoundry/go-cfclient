package client_test

import (
	"os"
	"path"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/client"
	"github.com/stretchr/testify/require"
)

func TestConfigNewUserPasswordConfig(t *testing.T) {
	c, err := client.NewUserPasswordConfig("https://api.example.com", "admin", "pwd")
	require.NoError(t, err)

	require.Equal(t, "Go-CF-client/2.0", c.UserAgent)

	require.Equal(t, "admin", c.Username)
	require.Equal(t, "pwd", c.Password)
	require.Equal(t, "https://api.example.com", c.ApiAddress)

	require.Empty(t, c.ClientID)
	require.Empty(t, c.ClientSecret)
}

func TestConfigNewUserPasswordConfigTrimsApiTrailingSlash(t *testing.T) {
	c, err := client.NewUserPasswordConfig("https://api.example.com/", "admin", "pwd")
	require.NoError(t, err)
	require.Equal(t, "https://api.example.com", c.ApiAddress)
}

func TestConfigNewUserPasswordConfigBadApiAddress(t *testing.T) {
	_, err := client.NewUserPasswordConfig("api.example.com", "admin", "pwd")
	require.Error(t, err)

	_, err = client.NewUserPasswordConfig("1.1.1.1", "admin", "pwd")
	require.Error(t, err)

	_, err = client.NewUserPasswordConfig("", "admin", "pwd")
	require.Error(t, err)
}

func TestConfigNewUserPasswordConfigEmptyUserPwd(t *testing.T) {
	_, err := client.NewUserPasswordConfig("https://api.example.com", "", "pwd")
	require.Error(t, err, "expected missing username error")

	_, err = client.NewUserPasswordConfig("https://api.example.com", "admin", "")
	require.Error(t, err, "expected missing password error")
}

func TestConfigNewClientSecretConfig(t *testing.T) {
	c, err := client.NewClientSecretConfig("https://api.example.com", "opsman", "secret")
	require.NoError(t, err)

	require.Equal(t, "Go-CF-client/2.0", c.UserAgent)

	require.Equal(t, "opsman", c.ClientID)
	require.Equal(t, "secret", c.ClientSecret)
	require.Equal(t, "https://api.example.com", c.ApiAddress)

	require.Empty(t, c.Username)
	require.Empty(t, c.Password)
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

	cfg, err := client.NewConfigFromCFHomeDir(cfHomeDir)
	require.NoError(t, err)

	require.Equal(t, "https://api.sys.example.com", cfg.ApiAddress)
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
  "AccessToken": "bearer secret-bearer-token",
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
