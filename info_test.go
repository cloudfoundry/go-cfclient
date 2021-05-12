package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetInfo(t *testing.T) {
	Convey("Get info", t, func() {
		setupMultiple(nil, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		info, err := client.GetInfo()
		So(err, ShouldBeNil)

		So(info.MinCLIVersion, ShouldEqual, "6.23.0")
	})

	Convey("Doesn't support metadata api", t, func() {
		setupMultiple([]MockRoute{
			{
				Method:   "GET",
				Endpoint: "/",
				Status:   200,
				Output: []string{`{
				   "links": {
				      "cloud_controller_v3": {
				         "href": "https://api.dev.cfdev.sh/v3",
				         "meta": {
				            "version": "3.65.0"
				         }
				      }
				   }
				}`},
			},
		}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		supports, err := client.SupportsMetadataAPI()
		So(err, ShouldBeNil)

		So(supports, ShouldEqual, false)
	})

	Convey("Support metadata api", t, func() {
		setupMultiple([]MockRoute{
			{
				Method:   "GET",
				Endpoint: "/",
				Status:   200,
				Output: []string{`{
				   "links": {
				      "cloud_controller_v3": {
				         "href": "https://api.dev.cfdev.sh/v3",
				         "meta": {
				            "version": "3.66.0"
				         }
				      }
				   }
				}`},
			},
		}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		supports, err := client.SupportsMetadataAPI()
		So(err, ShouldBeNil)

		So(supports, ShouldEqual, true)
	})
}
