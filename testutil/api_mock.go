package testutil

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	mux           *http.ServeMux
	server        *httptest.Server
	fakeUAAServer *httptest.Server
)

type MockRoute struct {
	Method           string
	Endpoint         string
	Output           []string
	UserAgent        string
	Status           int
	QueryString      string
	PostForm         string
	RedirectLocation string
}

func Setup(mock MockRoute, t *testing.T) string {
	return SetupMultiple([]MockRoute{mock}, t)
}

func SetupMultiple(mockEndpoints []MockRoute, t *testing.T) string {
	SetupFakeUAAServer(3)

	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	m := martini.New()
	m.Use(render.Renderer())
	r := martini.NewRouter()
	for _, mock := range mockEndpoints {
		method := mock.Method
		endpoint := mock.Endpoint
		output := mock.Output
		if len(output) == 0 {
			output = []string{""}
		}
		userAgent := mock.UserAgent
		status := mock.Status
		queryString := mock.QueryString
		postFormBody := mock.PostForm
		redirectLocation := mock.RedirectLocation
		switch method {
		case "GET":
			count := 0
			r.Get(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testQueryString(req.URL.RawQuery, queryString, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				singleOutput := output[count]
				count++
				return status, singleOutput
			})
		case "POST":
			r.Post(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testQueryString(req.URL.RawQuery, queryString, t)
				testReqBody(req, postFormBody, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				return status, output[0]
			})
		case "DELETE":
			r.Delete(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testQueryString(req.URL.RawQuery, queryString, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				return status, output[0]
			})
		case "PUT":
			r.Put(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testQueryString(req.URL.RawQuery, queryString, t)
				testReqBody(req, postFormBody, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				return status, output[0]
			})
		case "PATCH":
			r.Patch(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testQueryString(req.URL.RawQuery, queryString, t)
				testReqBody(req, postFormBody, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				return status, output[0]
			})
		case "PUT-FILE":
			r.Put(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testBodyContains(req, postFormBody, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				return status, output[0]
			})
		}
	}
	r.Get("/v2/info", func(r render.Render) {
		r.JSON(200, map[string]interface{}{
			"authorization_endpoint":       fakeUAAServer.URL,
			"token_endpoint":               fakeUAAServer.URL,
			"logging_endpoint":             server.URL,
			"name":                         "",
			"build":                        "",
			"support":                      "https://support.example.net",
			"version":                      0,
			"description":                  "",
			"min_cli_version":              "6.23.0",
			"min_recommended_cli_version":  "6.23.0",
			"api_version":                  "2.103.0",
			"app_ssh_endpoint":             "ssh.example.net:2222",
			"app_ssh_host_key_fingerprint": "00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:01",
			"app_ssh_oauth_client":         "ssh-proxy",
			"doppler_logging_endpoint":     "wss://doppler.example.net:443",
			"routing_endpoint":             "https://api.example.net/routing",
		})

	})

	m.Action(r.Handle)
	mux.Handle("/", m)

	return server.URL
}

func SetupFakeUAAServer(expiresIn int) {
	uaaMux := http.NewServeMux()
	fakeUAAServer = httptest.NewServer(uaaMux)
	m := martini.New()
	m.Use(render.Renderer())
	r := martini.NewRouter()
	count := 1
	r.Post("/oauth/token", func(r render.Render) {
		r.JSON(200, map[string]interface{}{
			"token_type":    "bearer",
			"access_token":  "foobar" + strconv.Itoa(count),
			"refresh_token": "barfoo",
			"expires_in":    expiresIn,
		})
		count = count + 1
	})
	r.NotFound(func() string { return "" })
	m.Action(r.Handle)
	uaaMux.Handle("/", m)
}

func Teardown() {
	server.Close()
	fakeUAAServer.Close()
}

func testQueryString(QueryString string, QueryStringExp string, t *testing.T) {
	t.Helper()
	if QueryStringExp == "" {
		return
	}

	value, _ := url.QueryUnescape(QueryString)
	if QueryStringExp != value {
		t.Errorf("Error: Query string '%s' should be equal to '%s'", QueryStringExp, value)
	}
}

func testUserAgent(UserAgent string, UserAgentExp string, t *testing.T) {
	t.Helper()
	if len(UserAgentExp) < 1 {
		UserAgentExp = "Go-CF-client/2.0"
	}
	if UserAgent != UserAgentExp {
		t.Errorf("Error: Agent %s should be equal to %s", UserAgent, UserAgentExp)
	}
}

func testReqBody(req *http.Request, postFormBody string, t *testing.T) {
	t.Helper()
	if postFormBody != "" {
		if body, err := io.ReadAll(req.Body); err != nil {
			t.Error("No request body but expected one")
		} else {
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(req.Body)
			require.JSONEq(t, postFormBody, string(body),
				"Expected request body (%s) does not equal request body (%s)", postFormBody, body)
		}
	}
}

func testBodyContains(req *http.Request, expected string, t *testing.T) {
	t.Helper()
	if expected != "" {
		if body, err := io.ReadAll(req.Body); err != nil {
			t.Error("No request body but expected one")
		} else {
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(req.Body)
			if !strings.Contains(string(body), expected) {
				t.Errorf("Expected request body (%s) was not found in actual request body (%s)", expected, body)
			}
		}
	}
}
