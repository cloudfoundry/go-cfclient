package cfclient

import (
	"github.com/go-martini/martini"
	"net/http"
	"net/http/httptest"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
)

type MockRoute struct {
	Method   string
	Endpoint string
	Output   string
}

func setup(mock MockRoute) {
	setupMultiple([]MockRoute{mock})
}

func setupMultiple(mockEndpoints []MockRoute) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	m := martini.New()
	r := martini.NewRouter()
	for _, mock := range mockEndpoints {
		method := mock.Method
		endpoint := mock.Endpoint
		output := mock.Output
		if method == "GET" {
			r.Get(endpoint, func() string {
				return output
			})
		} else if method == "POST" {
			r.Post(endpoint, func() string {
				return output
			})
		}
	}
	m.Action(r.Handle)
	mux.Handle("/", m)
}

func teardown() {
	server.Close()
}
