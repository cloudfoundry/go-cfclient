package cfclient

import (
	"testing"

	"github.com/onsi/gomega"
	. "github.com/smartystreets/goconvey/convey"
)

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

func TestMakeRequest(t *testing.T) {
	Convey("Test making request b", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", listOrgsPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress:        server.URL,
			Username:          "foo",
			Password:          "bar",
			SkipSslValidation: true,
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
		req := client.NewRequest("GET", "/v2/foobar")
		resp, err := client.DoRequest(req)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
	})
}

func TestTokenRefresh(t *testing.T) {
	gomega.RegisterTestingT(t)
	Convey("Test making request", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", listOrgsPayload, ""}, t)
		c := &Config{
			ApiAddress: server.URL,
			Username:   "foo",
			Password:   "bar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		token, err := client.GetToken()
		So(err, ShouldBeNil)

		gomega.Consistently(token).Should(gomega.Equal("bearer foobar2"))
		// gomega.Eventually(client.GetToken(), "3s").Should(gomega.Equal("bearer foobar3"))
	})
}
