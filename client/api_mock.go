package client

import (
	"encoding/json"
	"fmt"
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

type RouteTest struct {
	Description string
	Route       MockRoute
	Expected    string
	Expected2   string
	Expected3   string
	Action      func(c *Client, t *testing.T) (any, error)
	Action2     func(c *Client, t *testing.T) (any, any, error)
	Action3     func(c *Client, t *testing.T) (any, any, any, error)
}

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

func setup(mock MockRoute, t *testing.T) {
	setupMultiple([]MockRoute{mock}, t)
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

func setupMultiple(mockEndpoints []MockRoute, t *testing.T) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	fakeUAAServer = FakeUAAServer(3)
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
}

func FakeUAAServer(expiresIn int) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
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
	mux.Handle("/", m)
	return server
}

func teardown() {
	server.Close()
	fakeUAAServer.Close()
}

func executeTests(tests []RouteTest, t *testing.T) {
	for _, tt := range tests {
		func() {
			setup(tt.Route, t)
			defer teardown()
			details := fmt.Sprintf("%s %s", tt.Route.Method, tt.Route.Endpoint)
			if tt.Description != "" {
				details = tt.Description + ": " + details
			}

			c, _ := NewTokenConfig(server.URL, "foobar")
			cl, err := New(c)
			require.NoError(t, err, details)

			assertEq := func(t *testing.T, expected string, obj any) {
				if isJSON(expected) {
					actualJSON, err := json.Marshal(obj)
					require.NoError(t, err, details)
					require.JSONEq(t, expected, string(actualJSON), details)
				} else {
					if s, ok := obj.(string); ok {
						require.Equal(t, expected, s, details)
					}
				}
			}

			if tt.Action != nil {
				obj1, err := tt.Action(cl, t)
				require.NoError(t, err, details)
				assertEq(t, tt.Expected, obj1)
			} else if tt.Action2 != nil {
				obj1, obj2, err := tt.Action2(cl, t)
				require.NoError(t, err, details)
				assertEq(t, tt.Expected, obj1)
				assertEq(t, tt.Expected2, obj2)
			} else if tt.Action3 != nil {
				obj1, obj2, obj3, err := tt.Action3(cl, t)
				require.NoError(t, err, details)
				assertEq(t, tt.Expected, obj1)
				assertEq(t, tt.Expected2, obj2)
				assertEq(t, tt.Expected3, obj3)
			}
		}()
	}
}

func isJSON(obj string) bool {
	return strings.HasPrefix(obj, "{") || strings.HasPrefix(obj, "[")
}