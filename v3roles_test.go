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
			{"GET", "/v3/roles", []string{listV3SpaceRolesBySpaceGUIDPayload}, "", http.StatusOK, "space_guids=spaceGUID1", nil},
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

func TestListV3SpaceRolesByGUIDAndType(t *testing.T) {
	Convey("List V3 Space Roles By Space GUID and type", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/roles", []string{listV3SpaceRolesBySpaceGUIDAndTypePayload}, "", http.StatusOK, "include=user&space_guids=spaceGUID1&types=space_supporter", nil},
			{"GET", "/v3/rolespage2", []string{listV3SpaceRolesBySpaceGuidAndTypePayloadPage2}, "", http.StatusOK, "page=2&per_page=2&include=user&space_guids=spaceGUID1&types=space_supporter", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		users, err := client.ListV3SpaceRolesByGUIDAndType("spaceGUID1", "space_supporter")
		So(err, ShouldBeNil)
		So(users, ShouldHaveLength, 3)

		So(users[0].Username, ShouldEqual, "user1")
		So(users[1].Username, ShouldEqual, "user2")
		So(users[2].Username, ShouldEqual, "user3")
	})
}

func TestListV3OrgRolesByGUIDAndType(t *testing.T) {
	Convey("List V3 Org Roles By Org GUID and type", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/roles", []string{listV3OrganizationRolesByOrganizationGUIDAndTypePayload}, "", http.StatusOK, "include=user&organization_guids=orgGUID1&types=organization_auditor", nil},
			{"GET", "/v3/rolespage2", []string{listV3OrganizationRolesByOrganizationGuidAndTypePayloadPage2}, "", http.StatusOK, "page=2&per_page=2&include=user&organization_guids=orgGUID1&types=organization_auditor", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		users, err := client.ListV3OrganizationRolesByGUIDAndType("orgGUID1", "organization_auditor")
		So(err, ShouldBeNil)
		So(users, ShouldHaveLength, 3)

		So(users[0].Username, ShouldEqual, "user1")
		So(users[1].Username, ShouldEqual, "user2")
		So(users[2].Username, ShouldEqual, "user3")
	})
}

func TestCreateV3SpaceRoles(t *testing.T) {
	Convey("Create V3 Space Role", t, func() {
		setup(MockRoute{"POST", "/v3/roles", []string{createV3SpaceRolePayload}, "", http.StatusCreated, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		role, err := client.CreateV3SpaceRole(
			"b40a40c8-58b7-49a0-b47d-9d6fe5d72905",
			"c4958204-6b65-43ea-832b-e4c57aea6641",
			"space_supporter",
		)
		So(err, ShouldBeNil)
		So(role.Type, ShouldEqual, "space_supporter")
		So(role.GUID, ShouldEqual, "b9f59ab2-2b09-438e-bebb-30e8704ffb89")
		So(role.Links["self"].Href, ShouldEqual, "https://api.example.org/v3/roles/b9f59ab2-2b09-438e-bebb-30e8704ffb89")
		So(role.Links["user"].Href, ShouldEqual, "https://api.example.org/v3/users/c4958204-6b65-43ea-832b-e4c57aea6641")
		So(role.Links["space"].Href, ShouldEqual, "https://api.example.org/v3/spaces/b40a40c8-58b7-49a0-b47d-9d6fe5d72905")
	})
}

func TestCreateV3OrgRoles(t *testing.T) {
	Convey("Create V3 Org Role", t, func() {
		setup(MockRoute{"POST", "/v3/roles", []string{createV3OrganizationRolePayload}, "", http.StatusCreated, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		role, err := client.CreateV3OrganizationRole(
			"fa8a8346-0d92-4729-870c-77ee1934f973",
			"ac2e02c9-2c5c-4712-a620-a68449d263c3",
			"organization_user",
		)
		So(err, ShouldBeNil)
		So(role.Type, ShouldEqual, "organization_user")
		So(role.GUID, ShouldEqual, "21cbfaeb-bff7-4cfd-a7a9-6c13ec76f246")
		So(role.Links["self"].Href, ShouldEqual, "https://api.example.org/v3/roles/21cbfaeb-bff7-4cfd-a7a9-6c13ec76f246")
		So(role.Links["user"].Href, ShouldEqual, "https://api.example.org/v3/users/ac2e02c9-2c5c-4712-a620-a68449d263c3")
		So(role.Links["organization"].Href, ShouldEqual, "https://api.example.org/v3/organizations/fa8a8346-0d92-4729-870c-77ee1934f973")
	})
}

func TestDeleteV3Role(t *testing.T) {
	Convey("Delete V3 Role", t, func() {
		setup(MockRoute{"DELETE", "/v3/roles/1cb006ee-fb05-47e1-b541-c34179ddc446", []string{""}, "", http.StatusAccepted, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteV3Role("1cb006ee-fb05-47e1-b541-c34179ddc446")
		So(err, ShouldBeNil)
	})
}
