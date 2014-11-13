package cfclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-martini/martini"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultConfig(t *testing.T) {
	Convey("Default config", t, func() {
		c := DefaultConfig()
		So(c.ApiAddress, ShouldEqual, "https://api.10.244.0.34.xip.io")
		So(c.LoginAddress, ShouldEqual, "https://login.10.244.0.34.xip.io")
		So(c.Username, ShouldEqual, "admin")
		So(c.Password, ShouldEqual, "admin")
		So(c.Token, ShouldEqual, "")
	})
}

func TestCreateNewClient(t *testing.T) {
	Convey("Create new client", t, func() {
		c := &Config{
			ApiAddress:   "",
			LoginAddress: "",
			Username:     "",
			Password:     "",
		}
		client := NewClient(c)
		So(client, ShouldNotBeNil)
	})
}

func TestMakeRequest(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	m := martini.New()
	r := martini.NewRouter()
	r.Post("/oauth/token", func() string {
		return `{ "access_token": "foobar", "token_type": "test", "refresh_token": "blah"}`
	})
	r.NotFound(func() string { return "" })
	m.Action(r.Handle)
	mux.Handle("/", m)
	Convey("Test making request", t, func() {
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: server.URL,
			Username:     "foo",
			Password:     "bar",
		}
		client := NewClient(c)
		req := client.newRequest("GET", "/v2/foobar")
		resp, err := client.doRequest(req)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
	})
}
