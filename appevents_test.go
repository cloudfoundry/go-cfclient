package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListAppEvents(t *testing.T) {
	Convey("List App Events", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/events", listAppsCreatedEventPayload, ""},
			{"GET", "/v2/events2", listAppsCreatedEventPayload2, ""},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
		appEvents, err := client.ListAppEvents("blub")
		So(err.Error(), ShouldEqual, "Unsupported app event type blub")
		appEvents, err = client.ListAppEvents(AppCreate)
		So(err, ShouldEqual, nil)
		So(len(appEvents), ShouldEqual, 2)
		So(appEvents[0].MetaData.Request.State, ShouldEqual, "STOPPED")
		So(appEvents[1].MetaData.Request.State, ShouldEqual, "STARTED")
	})
}

func TestListAppEventsByQuery(t *testing.T) {
	Convey("List App Events By Query", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/events", listAppsCreatedEventPayload, ""},
			{"GET", "/v2/events2", listAppsCreatedEventPayload2, ""},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		appEvents, err := client.ListAppEventsByQuery("blub", []AppEventQuery{})
		So(err.Error(), ShouldEqual, "Unsupported app event type blub")

		appEventQuery := AppEventQuery{
			Filter:   "nofilter",
			Operator: ":",
			Value:    "retlifon",
		}
		appEvents, err = client.ListAppEventsByQuery(AppCreate, []AppEventQuery{appEventQuery})
		So(err.Error(), ShouldEqual, "Unsupported query filter type nofilter")

		appEventQuery = AppEventQuery{
			Filter:   FilterTimestamp,
			Operator: "not",
			Value:    "retlifon",
		}
		appEvents, err = client.ListAppEventsByQuery(AppCreate, []AppEventQuery{appEventQuery})
		So(err.Error(), ShouldEqual, "Unsupported query operator type not")

		appEventQuery = AppEventQuery{
			Filter:   FilterActee,
			Operator: ":",
			Value:    "3ca436ff-67a8-468a-8c7d-27ec68a6cfe5",
		}
		appEvents, err = client.ListAppEventsByQuery(AppCreate, []AppEventQuery{appEventQuery})
		So(err, ShouldEqual, nil)
		So(len(appEvents), ShouldEqual, 2)
		So(appEvents[0].MetaData.Request.State, ShouldEqual, "STOPPED")
		So(appEvents[1].MetaData.Request.State, ShouldEqual, "STARTED")
	})
}
