package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListSpaces(t *testing.T) {
	Convey("List Space", t, func() {
		setup(MockRoute{"GET", "/v2/spaces", listSpacesPayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		spaces := client.ListSpaces()
		So(len(spaces), ShouldEqual, 2)
		So(spaces[0].Guid, ShouldEqual, "8efd7c5c-d83c-4786-b399-b7bd548839e1")
		So(spaces[0].Name, ShouldEqual, "dev")
	})
}

func TestSpaceOrg(t *testing.T) {
	Convey("Find space org", t, func() {
		setup(MockRoute{"GET", "/v2/org/foobar", orgPayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		space := &Space{
			Guid:   "123",
			Name:   "test space",
			OrgURL: "/v2/org/foobar",
			c:      client,
		}
		org := space.Org()
		So(org.Name, ShouldEqual, "test-org")
		So(org.Guid, ShouldEqual, "da0dba14-6064-4f7a-b15a-ff9e677e49b2")
	})
}
