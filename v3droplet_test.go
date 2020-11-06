package cfclient

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetCurrentDropletForV3App(t *testing.T) {
	Convey("Set Droplet for V3 App", t, func() {
		body := `{"data":{"guid":"3fc0916f-2cea-4f3a-ae53-048388baa6bd"}}`
		setup(MockRoute{"PATCH", "/v3/apps/9d8e007c-ce52-4ea7-8a57-f2825d2c6b39/relationships/current_droplet", []string{currentDropletV3AppPayload}, "", http.StatusOK, "", &body}, t)
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
		setup(MockRoute{"GET", "/v3/apps/7b34f1cf-7e73-428a-bb5a-8a17a8058396/droplets/current", []string{getV3CurrentAppDropletPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.GetCurrentDropletForV3App("7b34f1cf-7e73-428a-bb5a-8a17a8058396")
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp.GUID, ShouldEqual, "585bc3c1-3743-497d-88b0-403ad6b56d16")
		So(resp.Links["self"].Href, ShouldEqual, "https://api.example.org/v3/droplets/585bc3c1-3743-497d-88b0-403ad6b56d16")
		So(resp.Links["assign_current_droplet"].Href, ShouldEqual, "https://api.example.org/v3/apps/7b34f1cf-7e73-428a-bb5a-8a17a8058396/relationships/current_droplet")
		So(resp.Links["assign_current_droplet"].Method, ShouldEqual, "PATCH")
	})
}

func TestDeleteDroplet(t *testing.T) {
	Convey("Delete Droplet", t, func() {
		setup(MockRoute{"DELETE", "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4", []string{""}, "", http.StatusAccepted, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteDroplet("59c3d133-2b83-46f3-960e-7765a129aea4")
		So(err, ShouldBeNil)
	})
}
