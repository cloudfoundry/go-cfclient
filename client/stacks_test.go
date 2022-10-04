package client

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListStacksByQuery(t *testing.T) {
	Convey("List All  stacks", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/stacks", []string{listStacksPayload}, "", http.StatusOK, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c, _ := NewTokenConfig(server.URL, "foobar")
		client, err := New(c)
		So(err, ShouldBeNil)

		stacks, err := client.ListStacksByQuery(nil)
		So(err, ShouldBeNil)
		So(stacks, ShouldHaveLength, 2)

		So(stacks[0].Name, ShouldEqual, "my-stack-1")
		So(stacks[0].Description, ShouldEqual, "This is my first stack!")
		So(stacks[0].GUID, ShouldEqual, "guid-1")
		So(stacks[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/stacks/guid-1")
		So(stacks[1].Name, ShouldEqual, "my-stack-2")
		So(stacks[1].Description, ShouldEqual, "This is my second stack!")
		So(stacks[1].GUID, ShouldEqual, "guid-2")
		So(stacks[1].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/stacks/guid-2")
	})
}
