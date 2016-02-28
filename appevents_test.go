package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListAppCreateEvent(t *testing.T) {
	Convey("List App Create Events", t, func() {
		//TODO enable tests with parameterized URLs
		//setup(MockRoute{"GET", "/v2/events?q=type:audit.app.create", listAppsCreatedEventPayload})
		setup(MockRoute{"GET", "/v2/events", listAppsCreatedEventPayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		orgs, err := client.ListAppCreateEvent()
		So(err, ShouldEqual, nil)
		So(len(orgs), ShouldEqual, 2)
	})
}
