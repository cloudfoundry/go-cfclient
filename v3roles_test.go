package cfclient

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListV3SpaceRolesByQuery(t *testing.T) {
	Convey("List V3 Space Roles By Space GUID", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/roles", []string{listV3SpaceRolesBySpaceGuidPayload}, "", http.StatusOK, "space_guids=spaceGUID1", nil},
			{"GET", "/v3/rolespage2", []string{listV3SpaceRolesBySpaceGuidPayloadPage2}, "", http.StatusOK, "page=2&per_page=2&space_guids=spaceGUID1", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		query := url.Values{}
		query["space_guids"] = []string{"spaceGUID1"}

		roles, err := client.ListV3RolesByQuery(query)
		So(err, ShouldBeNil)
		So(roles, ShouldHaveLength, 3)

		So(roles[0].Type, ShouldEqual, "space_developer")
		So(roles[1].Type, ShouldEqual, "space_auditor")
		So(roles[2].Type, ShouldEqual, "space_manager")

		So(roles[0].GUID, ShouldEqual, "roleGUID1")
		So(roles[0].Relationships["user"].Data.GUID, ShouldEqual, "userGUID1")
		So(roles[0].Relationships["space"].Data.GUID, ShouldEqual, "spaceGUID1")
		So(roles[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/roles/roleGUID1")
		So(roles[1].GUID, ShouldEqual, "roleGUID2")
		So(roles[1].Relationships["user"].Data.GUID, ShouldEqual, "userGUID2")
		So(roles[1].Relationships["space"].Data.GUID, ShouldEqual, "spaceGUID1")
		So(roles[1].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/roles/roleGUID2")
		So(roles[2].GUID, ShouldEqual, "roleGUID3")
		So(roles[2].Relationships["user"].Data.GUID, ShouldEqual, "userGUID2")
		So(roles[2].Relationships["space"].Data.GUID, ShouldEqual, "spaceGUID1")
		So(roles[2].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/roles/roleGUID3")
	})

	Convey("List V3 Space Roles By User GUID", t, func() {
		setup(MockRoute{"GET", "/v3/roles", []string{listV3SpaceRolesByUserGuidPayload}, "", http.StatusOK, "user_guids=userGUID1", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		query := url.Values{}
		query["user_guids"] = []string{"userGUID1"}

		roles, err := client.ListV3RolesByQuery(query)
		So(err, ShouldBeNil)
		So(roles, ShouldHaveLength, 2)

		So(roles[0].Type, ShouldEqual, "space_developer")
		So(roles[1].Type, ShouldEqual, "space_manager")

		So(roles[0].GUID, ShouldEqual, "roleGUID1")
		So(roles[0].Relationships["user"].Data.GUID, ShouldEqual, "userGUID1")
		So(roles[0].Relationships["space"].Data.GUID, ShouldEqual, "spaceGUID1")
		So(roles[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/roles/roleGUID1")
		So(roles[1].GUID, ShouldEqual, "roleGUID4")
		So(roles[1].Relationships["user"].Data.GUID, ShouldEqual, "userGUID1")
		So(roles[1].Relationships["space"].Data.GUID, ShouldEqual, "spaceGUID2")
		So(roles[1].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/roles/roleGUID4")
	})

	Convey("List V3 Space Users By Space Guid and User GUID", t, func() {
		setup(MockRoute{"GET", "/v3/roles", []string{listV3spaceRolesBySpaceAndUserGuidPayload}, "", http.StatusOK, "space_guids=spaceGUID2&user_guids=userGUID1", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		query := url.Values{}
		query["space_guids"] = []string{"spaceGUID2"}
		query["user_guids"] = []string{"userGUID1"}

		roles, err := client.ListV3RolesByQuery(query)
		So(err, ShouldBeNil)
		So(roles, ShouldHaveLength, 1)

		So(roles[0].Type, ShouldEqual, "space_manager")

		So(roles[0].GUID, ShouldEqual, "roleGUID4")
		So(roles[0].Relationships["user"].Data.GUID, ShouldEqual, "userGUID1")
		So(roles[0].Relationships["space"].Data.GUID, ShouldEqual, "spaceGUID2")
		So(roles[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/roles/roleGUID4")
	})
}
