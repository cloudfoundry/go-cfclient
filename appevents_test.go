package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListAppCreateEvents(t *testing.T) {
	Convey("List App Create Events", t, func() {
		setup(MockRoute{"GET", "/v2/events", listAppsCreatedEventPayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		appEvents, err := client.ListAppEvents("blub")
		So(err.Error(), ShouldEqual, "Unsupported app event type blub")
		appEvents, err = client.ListAppEvents(AppCreate)
		So(err, ShouldEqual, nil)
		So(len(appEvents), ShouldEqual, 2)

	})
}
