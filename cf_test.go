package cfclient

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-martini/martini"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
)

func setup(method, endpoint, output string) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	m := martini.New()
	r := martini.NewRouter()
	if method == "GET" {
		r.Get(endpoint, func() string {
			return output
		})
	} else if method == "POST" {
		r.Post(endpoint, func() string {
			return output
		})
	}
	m.Action(r.Handle)
	mux.Handle("/", m)
}

func teardown() {
	server.Close()
}
