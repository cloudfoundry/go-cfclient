package cfclient

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

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

func TestNewConfigFromCFHome(t *testing.T) {
	Convey("New config from CF_HOME", t, func() {
		cfHomeDir, err := ioutil.TempDir("", "cf_home")
		So(err, ShouldBeNil)

		configDir := path.Join(cfHomeDir, ".cf")
		err = os.MkdirAll(configDir, 0744)
		So(err, ShouldBeNil)

		configPath := path.Join(configDir, "config.json")
		err = os.WriteFile(configPath, []byte(cfCLIConfig), 0744)
		So(err, ShouldBeNil)

		cfg, err := NewConfigFromCFHome(cfHomeDir)
		So(err, ShouldBeNil)

		So(cfg.Token, ShouldEqual, "secret-bearer-token")
		So(cfg.ApiAddress, ShouldEqual, "https://api.sys.example.com")
		So(cfg.SkipSslValidation, ShouldBeTrue)
	})
}

func TestDefaultConfig(t *testing.T) {
	Convey("Default config", t, func() {
		c := DefaultConfig()
		So(c.ApiAddress, ShouldEqual, "http://api.bosh-lite.com")
		So(c.Username, ShouldEqual, "admin")
		So(c.Password, ShouldEqual, "admin")
		So(c.SkipSslValidation, ShouldEqual, false)
		So(c.Token, ShouldEqual, "")
		So(c.UserAgent, ShouldEqual, "Go-CF-client/1.1")
	})
}

func TestRemovalofTrailingSlashOnAPIAddress(t *testing.T) {
	Convey("Test removal of trailing slash of the API Address", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", []string{listOrgsPayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL + "/",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
		So(client.Config.ApiAddress, ShouldNotEndWith, "/")
	})
}

func TestMakeRequest(t *testing.T) {
	Convey("Test making request b", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", []string{listOrgsPayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress:        server.URL,
			Username:          "foo",
			Password:          "bar",
			SkipSslValidation: true,
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
		req := client.NewRequest("GET", "/v2/organizations")
		resp, err := client.DoRequest(req)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
	})
}

func TestMakeRequestFailure(t *testing.T) {
	Convey("Test making request b", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", []string{listOrgsPayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress:        server.URL,
			Username:          "foo",
			Password:          "bar",
			SkipSslValidation: true,
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
		req := client.NewRequest("GET", "/v2/organizations")
		req.url = "%gh&%ij"
		resp, err := client.DoRequest(req)
		So(resp, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestMakeRequestWithTimeout(t *testing.T) {
	Convey("Test making request b", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", []string{listOrgsPayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress:        server.URL,
			Username:          "foo",
			Password:          "bar",
			SkipSslValidation: true,
			HttpClient:        &http.Client{Timeout: 10 * time.Nanosecond},
		}
		client, err := NewClient(c)
		So(err, ShouldNotBeNil)
		So(client, ShouldBeNil)
	})
}

func TestHTTPErrorHandling(t *testing.T) {
	Convey("Test making request b", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", []string{"502 Bad Gateway"}, "", 502, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress:        server.URL,
			Username:          "foo",
			Password:          "bar",
			SkipSslValidation: true,
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
		req := client.NewRequest("GET", "/v2/organizations")
		resp, err := client.DoRequest(req)
		So(err, ShouldNotBeNil)
		So(resp, ShouldNotBeNil)

		httpErr := err.(CloudFoundryHTTPError)
		So(httpErr.StatusCode, ShouldEqual, 502)
		So(httpErr.Status, ShouldEqual, "502 Bad Gateway")
		So(string(httpErr.Body), ShouldEqual, "502 Bad Gateway")
	})
}

func TestTokenRefresh(t *testing.T) {
	Convey("Test making request", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", []string{listOrgsPayload}, "", 200, "", nil}, t)
		fakeUAAServer = FakeUAAServer(1)
		c := &Config{
			ApiAddress: server.URL,
			Username:   "foo",
			Password:   "bar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		token, err := client.GetToken()
		So(err, ShouldBeNil)
		So(token, ShouldEqual, "bearer foobar2")

		for i := 0; i < 5; i++ {
			token, _ = client.GetToken()
			if token == "bearer foobar3" {
				break
			}
			time.Sleep(time.Second)
		}
		So(token, ShouldEqual, "bearer foobar3")
	})
}

func TestEndpointRefresh(t *testing.T) {
	Convey("Test expiring endpoint", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", []string{listOrgsPayload}, "", 200, "", nil}, t)
		fakeUAAServer = FakeUAAServer(0)

		c := &Config{
			ApiAddress: server.URL,
			Username:   "foo",
			Password:   "bar",
		}

		client, err := NewClient(c)
		So(err, ShouldBeNil)

		lastTokenSource := client.Config.TokenSource
		for i := 1; i < 5; i++ {
			_, err := client.GetToken()
			So(err, ShouldBeNil)
			So(client.Config.TokenSource, ShouldNotEqual, lastTokenSource)
			lastTokenSource = client.Config.TokenSource
		}
	})
}
