package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListServices(t *testing.T) {
	Convey("List Services", t, func() {
		setup(MockRoute{"GET", "/v2/services", listServicePayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		services := client.ListServices()
		So(len(services), ShouldEqual, 2)
		So(services[0].Guid, ShouldEqual, "a3d76c01-c08a-4505-b06d-8603265682a3")
		So(services[0].Label, ShouldEqual, "nats")
	})
}
