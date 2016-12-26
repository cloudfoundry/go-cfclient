package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListSpaces(t *testing.T) {
	Convey("List Space", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/spaces", listSpacesPayload, ""},
			{"GET", "/v2/spacesPage2", listSpacesPayloadPage2, ""},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		spaces, err := client.ListSpaces()
		So(err, ShouldBeNil)

		So(len(spaces), ShouldEqual, 4)
		So(spaces[0].Guid, ShouldEqual, "8efd7c5c-d83c-4786-b399-b7bd548839e1")
		So(spaces[0].Name, ShouldEqual, "dev")
		So(spaces[1].Guid, ShouldEqual, "657b5923-7de0-486a-9928-b4d78ee24931")
		So(spaces[1].Name, ShouldEqual, "demo")
		So(spaces[2].Guid, ShouldEqual, "9ffd7c5c-d83c-4786-b399-b7bd54883977")
		So(spaces[2].Name, ShouldEqual, "test")
		So(spaces[3].Guid, ShouldEqual, "329b5923-7de0-486a-9928-b4d78ee24982")
		So(spaces[3].Name, ShouldEqual, "prod")
	})
}

func TestSpaceOrg(t *testing.T) {
	Convey("Find space org", t, func() {
		setup(MockRoute{"GET", "/v2/org/foobar", orgPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space := &Space{
			Guid:   "123",
			Name:   "test space",
			OrgURL: "/v2/org/foobar",
			c:      client,
		}
		org, err := space.Org()
		So(err, ShouldBeNil)

		So(org.Name, ShouldEqual, "test-org")
		So(org.Guid, ShouldEqual, "da0dba14-6064-4f7a-b15a-ff9e677e49b2")
	})
}
