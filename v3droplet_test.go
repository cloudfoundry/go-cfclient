package cfclient

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetCurrentDropletForV3App(t *testing.T) {
	Convey("Set Droplet for V3 App", t, func() {
		body := `{"data":{"guid":"3fc0916f-2cea-4f3a-ae53-048388baa6bd"}}`
		setup(MockRoute{"PATCH", "/v3/apps/9d8e007c-ce52-4ea7-8a57-f2825d2c6b39/relationships/current_droplet", currentDropletV3AppPayload, "", http.StatusOK, "", &body}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.SetCurrentDropletForV3App("9d8e007c-ce52-4ea7-8a57-f2825d2c6b39", "3fc0916f-2cea-4f3a-ae53-048388baa6bd")
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp.Data.GUID, ShouldEqual, "9d8e007c-ce52-4ea7-8a57-f2825d2c6b39")
		So(resp.Links["self"].Href, ShouldEqual, "https://api.example.org/v3/apps/d4c91047-7b29-4fda-b7f9-04033e5c9c9f/relationships/current_droplet")
		So(resp.Links["related"].Href, ShouldEqual, "https://api.example.org/v3/apps/d4c91047-7b29-4fda-b7f9-04033e5c9c9f/droplets/current")
	})
}

func TestGetCurrentDropletForV3App(t *testing.T) {
	Convey("Get Droplet for V3 App", t, func() {
		setup(MockRoute{"GET", "/v3/apps/9d8e007c-ce52-4ea7-8a57-f2825d2c6b39/relationships/current_droplet", currentDropletV3AppPayload, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.GetCurrentDropletForV3App("9d8e007c-ce52-4ea7-8a57-f2825d2c6b39")
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp.Data.GUID, ShouldEqual, "9d8e007c-ce52-4ea7-8a57-f2825d2c6b39")
		So(resp.Links["self"].Href, ShouldEqual, "https://api.example.org/v3/apps/d4c91047-7b29-4fda-b7f9-04033e5c9c9f/relationships/current_droplet")
		So(resp.Links["related"].Href, ShouldEqual, "https://api.example.org/v3/apps/d4c91047-7b29-4fda-b7f9-04033e5c9c9f/droplets/current")
	})
}
