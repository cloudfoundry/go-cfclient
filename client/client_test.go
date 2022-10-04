package client

import (
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMakeRequest(t *testing.T) {
	Convey("Test making request b", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", []string{listOrgsPayload}, "", 200, "", nil}, t)
		defer teardown()
		c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
		c.SkipSSLValidation(true)
		client, err := New(c)
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
		c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
		c.SkipSSLValidation(true)
		client, err := New(c)
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
		c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
		c.SkipSSLValidation(true)
		c.HTTPClient(&http.Client{Timeout: 10 * time.Nanosecond})
		client, err := New(c)
		So(err, ShouldNotBeNil)
		So(client, ShouldBeNil)
	})
}

func TestHTTPErrorHandling(t *testing.T) {
	Convey("Test making request b", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", []string{"502 Bad Gateway"}, "", 502, "", nil}, t)
		defer teardown()
		c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
		c.SkipSSLValidation(true)
		client, err := New(c)
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
		c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
		client, err := New(c)
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
		c, _ := NewUserPasswordConfig(server.URL, "foo", "bar")
		client, err := New(c)
		So(err, ShouldBeNil)

		//lastTokenSource := client.Config.TokenSource
		for i := 1; i < 5; i++ {
			_, err := client.GetToken()
			So(err, ShouldBeNil)
			//So(client.Config.TokenSource, ShouldNotEqual, lastTokenSource)
			//lastTokenSource = client.Config.TokenSource
		}
	})
}
