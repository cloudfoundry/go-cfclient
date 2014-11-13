package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListSpaces(t *testing.T) {
	Convey("List Space", t, func() {
		setup("GET", "/v2/spaces", listSpacesPayload)
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: server.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		spaces := client.ListSpaces()
		So(len(spaces), ShouldEqual, 2)
		So(spaces[0].Guid, ShouldEqual, "8efd7c5c-d83c-4786-b399-b7bd548839e1")
		So(spaces[0].Name, ShouldEqual, "dev")
	})
}
