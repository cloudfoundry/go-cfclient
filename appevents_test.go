package cfclient

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListAppCreateEvent(t *testing.T) {
	Convey("List App Create Events", t, func() {
		//TODO enable tests with parameterized URLs
		setup(MockRoute{"GET", "/v2/events", listAppsCreatedEventPayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		appEvents, err := client.ListAppCreateEvent()
		fmt.Println(appEvents)
		So(err, ShouldEqual, nil)
		So(len(appEvents), ShouldEqual, 2)
	})
}

func TestListAppCreateEvent2(t *testing.T) {
	Convey("List App Create Events", t, func() {
		//TODO enable tests with parameterized URLs
		setup(MockRoute{"GET", "/v2/events", listOrgsPayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		appEvents, err := client.ListAppCreateEvent()
		fmt.Println(appEvents)
		So(err, ShouldEqual, nil)
		So(len(appEvents), ShouldEqual, 2)
	})
}
