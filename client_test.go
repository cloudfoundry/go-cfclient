package cfclient

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/onsi/gomega"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultConfig(t *testing.T) {
	Convey("Default config", t, func() {
		c := DefaultConfig()
		So(c.ApiAddress, ShouldEqual, "https://api.10.244.0.34.xip.io")
		So(c.LoginAddress, ShouldEqual, "https://login.10.244.0.34.xip.io")
		So(c.Username, ShouldEqual, "admin")
		So(c.Password, ShouldEqual, "admin")
		So(c.SkipSslValidation, ShouldEqual, false)
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
	Convey("Test making request", t, func() {
		server := fakeUAAServer()
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

func TestTokenRefresh(t *testing.T) {
	gomega.RegisterTestingT(t)
	Convey("Test making request", t, func() {
		server := fakeUAAServer()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: server.URL,
			Username:     "foo",
			Password:     "bar",
		}
		client := NewClient(c)
		gomega.Consistently(client.GetToken()).Should(gomega.Equal("bearer foobar2"))
		gomega.Eventually(client.GetToken(), "3s").Should(gomega.Equal("bearer foobar3"))
	})
}

func fakeUAAServer() *httptest.Server {
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
