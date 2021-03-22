package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetProcessStats(t *testing.T) {
	Convey("Get Process Stats", t, func() {
		setup(MockRoute{"GET", "/v3/processes/9902530c-c634-4864-a189-71d763cb12e2/stats", []string{getProcessStatsPayload1}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		stats, err := client.GetProcessStats("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(stats[0].State, ShouldEqual, "RUNNING")
		So(stats[0].Type, ShouldEqual, "web")

	})
}
