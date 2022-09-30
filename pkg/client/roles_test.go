package client

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListSpaceRolesByQuery(t *testing.T) {
	Convey("List  Space Roles By Space GUID", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/roles", []string{listSpaceRolesBySpaceGUIDPayload}, "", http.StatusOK, "space_guids=spaceGUID1", nil},
			{"GET", "/v3/rolespage2", []string{listSpaceRolesBySpaceGuidPayloadPage2}, "", http.StatusOK, "page=2&per_page=2&space_guids=spaceGUID1", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		query := url.Values{}
		query["space_guids"] = []string{"spaceGUID1"}

		roles, err := client.ListRolesByQuery(query)
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

	Convey("List  Space Roles By User GUID", t, func() {
		setup(MockRoute{"GET", "/v3/roles", []string{listSpaceRolesByUserGuidPayload}, "", http.StatusOK, "user_guids=userGUID1", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		query := url.Values{}
		query["user_guids"] = []string{"userGUID1"}

		roles, err := client.ListRolesByQuery(query)
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

	Convey("List  Space Users By Space Guid and User GUID", t, func() {
		setup(MockRoute{"GET", "/v3/roles", []string{listSpaceRolesBySpaceAndUserGuidPayload}, "", http.StatusOK, "space_guids=spaceGUID2&user_guids=userGUID1", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		query := url.Values{}
		query["space_guids"] = []string{"spaceGUID2"}
		query["user_guids"] = []string{"userGUID1"}

		roles, err := client.ListRolesByQuery(query)
		So(err, ShouldBeNil)
		So(roles, ShouldHaveLength, 1)

		So(roles[0].Type, ShouldEqual, "space_manager")

		So(roles[0].GUID, ShouldEqual, "roleGUID4")
		So(roles[0].Relationships["user"].Data.GUID, ShouldEqual, "userGUID1")
		So(roles[0].Relationships["space"].Data.GUID, ShouldEqual, "spaceGUID2")
		So(roles[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/roles/roleGUID4")
	})

}

func TestListSpaceRolesByGUID(t *testing.T) {
	Convey("List  Space Roles By Space GUID and type", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/roles", []string{listSpaceRoleUsersBySpaceGUIDPayload}, "", http.StatusOK, "include=user&space_guids=spaceGUID1", nil},
			{"GET", "/v3/rolespage2", []string{listSpaceRoleUsersBySpaceGUIDPayloadPage2}, "", http.StatusOK, "page=2&per_page=2&include=user&space_guids=spaceGUID1", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		roles, users, err := client.ListSpaceRolesByGUID("spaceGUID1")
		So(err, ShouldBeNil)
		So(roles, ShouldHaveLength, 3)
		So(users, ShouldHaveLength, 3)

		So(users[0].Username, ShouldEqual, "user1")
		So(users[1].Username, ShouldEqual, "user2")
		So(users[2].Username, ShouldEqual, "user3")
	})
}

func TestListSpaceRolesByGUIDAndType(t *testing.T) {
	Convey("List  Space Roles By Space GUID and type", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/roles", []string{listSpaceRolesBySpaceGUIDAndTypePayload}, "", http.StatusOK, "include=user&space_guids=spaceGUID1&types=space_supporter", nil},
			{"GET", "/v3/rolespage2", []string{listSpaceRolesBySpaceGuidAndTypePayloadPage2}, "", http.StatusOK, "page=2&per_page=2&include=user&space_guids=spaceGUID1&types=space_supporter", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		users, err := client.ListSpaceRolesByGUIDAndType("spaceGUID1", "space_supporter")
		So(err, ShouldBeNil)
		So(users, ShouldHaveLength, 3)

		So(users[0].Username, ShouldEqual, "user1")
		So(users[1].Username, ShouldEqual, "user2")
		So(users[2].Username, ShouldEqual, "user3")
	})
}

func TestListOrgRolesByGUID(t *testing.T) {
	Convey("List  Org Roles By Org GUID and type", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/roles", []string{listOrganizationRolesByOrganizationGUIDPayload}, "", http.StatusOK, "include=user&organization_guids=orgGUID1", nil},
			{"GET", "/v3/rolespage2", []string{listOrganizationRolesByOrganizationGuidPayloadPage2}, "", http.StatusOK, "page=2&per_page=2&include=user&organization_guids=orgGUID1", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		roles, users, err := client.ListOrganizationRolesByGUID("orgGUID1")
		So(err, ShouldBeNil)
		So(roles, ShouldHaveLength, 3)
		So(users, ShouldHaveLength, 3)

		So(users[0].Username, ShouldEqual, "user1")
		So(users[1].Username, ShouldEqual, "user2")
		So(users[2].Username, ShouldEqual, "user3")
	})
}

func TestListOrgRolesByGUIDAndType(t *testing.T) {
	Convey("List  Org Roles By Org GUID and type", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/roles", []string{listOrganizationRolesByOrganizationGUIDAndTypePayload}, "", http.StatusOK, "include=user&organization_guids=orgGUID1&types=organization_auditor", nil},
			{"GET", "/v3/rolespage2", []string{listOrganizationRolesByOrganizationGuidAndTypePayloadPage2}, "", http.StatusOK, "page=2&per_page=2&include=user&organization_guids=orgGUID1&types=organization_auditor", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		users, err := client.ListOrganizationRolesByGUIDAndType("orgGUID1", "organization_auditor")
		So(err, ShouldBeNil)
		So(users, ShouldHaveLength, 3)

		So(users[0].Username, ShouldEqual, "user1")
		So(users[1].Username, ShouldEqual, "user2")
		So(users[2].Username, ShouldEqual, "user3")
	})
}

func TestCreateSpaceRoles(t *testing.T) {
	Convey("Create  Space Role", t, func() {
		setup(MockRoute{"POST", "/v3/roles", []string{createSpaceRolePayload}, "", http.StatusCreated, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		role, err := client.CreateSpaceRole(
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

func TestCreateOrgRoles(t *testing.T) {
	Convey("Create  Org Role", t, func() {
		setup(MockRoute{"POST", "/v3/roles", []string{createOrganizationRolePayload}, "", http.StatusCreated, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		role, err := client.CreateOrganizationRole(
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

func TestDeleteRole(t *testing.T) {
	Convey("Delete  Role", t, func() {
		setup(MockRoute{"DELETE", "/v3/roles/1cb006ee-fb05-47e1-b541-c34179ddc446", []string{""}, "", http.StatusAccepted, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteRole("1cb006ee-fb05-47e1-b541-c34179ddc446")
		So(err, ShouldBeNil)
	})
}
