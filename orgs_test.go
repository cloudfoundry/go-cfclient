package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListOrgs(t *testing.T) {
	Convey("List Org", t, func() {
		setup(MockRoute{"GET", "/v2/organizations", listOrgsPayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		orgs := client.ListOrgs()
		So(len(orgs), ShouldEqual, 2)
		So(orgs[0].Guid, ShouldEqual, "a537761f-9d93-4b30-af17-3d73dbca181b")
		So(orgs[0].Name, ShouldEqual, "demo")
	})
}

func TestOrgSpaces(t *testing.T) {
	Convey("Get spaces by org", t, func() {
		setup(MockRoute{"GET", "/v2/organizations/foo/spaces", orgSpacesPayload})
		defer teardown()
		c := &Config{
			ApiAddress:   server.URL,
			LoginAddress: fakeUAAServer.URL,
			Token:        "foobar",
		}
		client := NewClient(c)
		spaces := client.OrgSpaces("foo")
		So(len(spaces), ShouldEqual, 1)
		So(spaces[0].Guid, ShouldEqual, "b8aff561-175d-45e8-b1e7-67e2aedb03b6")
		So(spaces[0].Name, ShouldEqual, "test")
	})
}
