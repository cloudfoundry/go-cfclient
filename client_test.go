package cfclient

import (
	"github.com/onsi/gomega"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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

func TestMakeRequest(t *testing.T) {
	Convey("Test making request", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", listOrgsPayload})
		defer teardown()
		c := &Config{
			ApiAddress:        server.URL,
			LoginAddress:      fakeUAAServer.URL,
			Username:          "foo",
			Password:          "bar",
			SkipSslValidation: true,
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
		setup(MockRoute{"GET", "/v2/organizations", listOrgsPayload})
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Username:     "foo",
			Password:     "bar",
		}
		client := NewClient(c)
		gomega.Consistently(client.GetToken()).Should(gomega.Equal("bearer foobar2"))
		gomega.Eventually(client.GetToken(), "3s").Should(gomega.Equal("bearer foobar3"))
	})
}
