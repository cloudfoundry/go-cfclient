package cfclient

import (
	"net/http"
	"net/http/httptest"
	"strconv"
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
	Method    string
	Endpoint  string
	Output    string
	UserAgent string
}

func setup(mock MockRoute, t *testing.T) {
	setupMultiple([]MockRoute{mock}, t)
}

func testUserAgent(UserAgent string, UserAgentExp string, t *testing.T) {
	if len(UserAgentExp) < 1 {
		UserAgentExp = "Go-CF-client/1.1"
	}
	if UserAgent != UserAgentExp {
		t.Fatalf("Error Agent %s should be equal to %s", UserAgent, UserAgentExp)
	}
}

func setupMultiple(mockEndpoints []MockRoute, t *testing.T) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	fakeUAAServer = FakeUAAServer()
	m := martini.New()
	m.Use(render.Renderer())
	r := martini.NewRouter()
	for _, mock := range mockEndpoints {
		method := mock.Method
		endpoint := mock.Endpoint
		output := mock.Output
		userAgent := mock.UserAgent
		if method == "GET" {
			r.Get(endpoint, func(req *http.Request) string {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				return output
			})
		} else if method == "POST" {
			r.Post(endpoint, func(req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				return 201, output
			})
		} else if method == "DELETE" {
			r.Delete(endpoint, func(req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				return 204, output
			})
		} else if method == "PUT" {
			r.Put(endpoint, func(req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				return 202, output
			})
		}
	}
	r.Get("/v2/info", func(r render.Render) {
		r.JSON(200, map[string]interface{}{
			"authorization_endpoint": fakeUAAServer.URL,
			"token_endpoint":         fakeUAAServer.URL,
			"logging_endpoint":       server.URL,
		})

	})

	m.Action(r.Handle)
	mux.Handle("/", m)
}

func FakeUAAServer() *httptest.Server {
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
			"expires_in":    3,
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
