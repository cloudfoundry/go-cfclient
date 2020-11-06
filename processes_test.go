package cfclient

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListProcesses(t *testing.T) {
	Convey("List Processes", t, func() {
		setup(MockRoute{"GET", "/v3/processes", []string{listProcessesPayload1, listProcessesPayload2}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		q := url.Values{}
		q.Add("per_page", "20")
		procs, err := client.ListAllProcessesByQuery(q)
		So(err, ShouldBeNil)

		So(procs, ShouldHaveLength, 26)
	})
}
