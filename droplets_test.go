package cfclient

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetDropletsCurrent(t *testing.T) {
	Convey("Get App Droplets Current", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/apps/faf9bc88-969d-4fdb-b2ee-d005ccb056cb/droplets/current", appCurrentDropletPayload, "Test-golang", 200, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
			UserAgent:  "Test-golang",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		droplet, err := client.GetAppsDropletsCurrent("faf9bc88-969d-4fdb-b2ee-d005ccb056cb")

		So(err, ShouldBeNil)

		So(droplet.Guid, ShouldEqual, "dd44a00f-3a37-4835-8d87-c1af920a8a39")
		So(droplet.CreatedAt.String(), ShouldEqual, "2020-09-23 18:14:18 +0000 UTC")
	})
}
