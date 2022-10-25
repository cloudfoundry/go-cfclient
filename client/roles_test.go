package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestListSpaceRolesByQueryWithSpaceGUID(t *testing.T) {
	mocks := []MockRoute{
		{"GET", "/v3/roles", []string{listSpaceRolesBySpaceGUIDPayload}, "", http.StatusOK, "space_guids=spaceGUID1", ""},
		{"GET", "/v3/rolespage2", []string{listSpaceRolesBySpaceGuidPayloadPage2}, "", http.StatusOK, "page=2&per_page=2&space_guids=spaceGUID1", ""},
	}
	setupMultiple(mocks, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	query := url.Values{}
	query["space_guids"] = []string{"spaceGUID1"}

	roles, err := client.Roles.ListRolesByQuery(query)
	require.NoError(t, err)
	require.Len(t, roles, 3)

	require.Equal(t, "space_developer", roles[0].Type)
	require.Equal(t, "space_auditor", roles[1].Type)
	require.Equal(t, "space_manager", roles[2].Type)

	require.Equal(t, "roleGUID1", roles[0].GUID)
	require.Equal(t, "userGUID1", roles[0].Relationships["user"].Data.GUID)
	require.Equal(t, "spaceGUID1", roles[0].Relationships["space"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/roles/roleGUID1", roles[0].Links["self"].Href)
	require.Equal(t, "roleGUID2", roles[1].GUID)
	require.Equal(t, "userGUID2", roles[1].Relationships["user"].Data.GUID)
	require.Equal(t, "spaceGUID1", roles[1].Relationships["space"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/roles/roleGUID2", roles[1].Links["self"].Href)
	require.Equal(t, "roleGUID3", roles[2].GUID)
	require.Equal(t, "userGUID2", roles[2].Relationships["user"].Data.GUID)
	require.Equal(t, "spaceGUID1", roles[2].Relationships["space"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/roles/roleGUID3", roles[2].Links["self"].Href)
}
func TestListSpaceRolesByQueryWithUserGUID(t *testing.T) {
	setup(MockRoute{"GET", "/v3/roles", []string{listSpaceRolesByUserGuidPayload}, "", http.StatusOK, "user_guids=userGUID1", ""}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	query := url.Values{}
	query["user_guids"] = []string{"userGUID1"}

	roles, err := client.Roles.ListRolesByQuery(query)
	require.NoError(t, err)
	require.Len(t, roles, 2)

	require.Equal(t, "space_developer", roles[0].Type)
	require.Equal(t, "space_manager", roles[1].Type)

	require.Equal(t, "roleGUID1", roles[0].GUID)
	require.Equal(t, "userGUID1", roles[0].Relationships["user"].Data.GUID)
	require.Equal(t, "spaceGUID1", roles[0].Relationships["space"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/roles/roleGUID1", roles[0].Links["self"].Href)
	require.Equal(t, "roleGUID4", roles[1].GUID)
	require.Equal(t, "userGUID1", roles[1].Relationships["user"].Data.GUID)
	require.Equal(t, "spaceGUID2", roles[1].Relationships["space"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/roles/roleGUID4", roles[1].Links["self"].Href)
}

func TestListSpaceRolesByQueryWithSpaceAndUserGUID(t *testing.T) {
	setup(MockRoute{"GET", "/v3/roles", []string{listSpaceRolesBySpaceAndUserGuidPayload}, "", http.StatusOK, "space_guids=spaceGUID2&user_guids=userGUID1", ""}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	query := url.Values{}
	query["space_guids"] = []string{"spaceGUID2"}
	query["user_guids"] = []string{"userGUID1"}

	roles, err := client.Roles.ListRolesByQuery(query)
	require.NoError(t, err)
	require.Len(t, roles, 1)

	require.Equal(t, "space_manager", roles[0].Type)

	require.Equal(t, "roleGUID4", roles[0].GUID)
	require.Equal(t, "userGUID1", roles[0].Relationships["user"].Data.GUID)
	require.Equal(t, "spaceGUID2", roles[0].Relationships["space"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/roles/roleGUID4", roles[0].Links["self"].Href)

}

func TestListSpaceRolesByGUID(t *testing.T) {
	mocks := []MockRoute{
		{"GET", "/v3/roles", []string{listSpaceRoleUsersBySpaceGUIDPayload}, "", http.StatusOK, "include=user&space_guids=spaceGUID1", ""},
		{"GET", "/v3/rolespage2", []string{listSpaceRoleUsersBySpaceGUIDPayloadPage2}, "", http.StatusOK, "page=2&per_page=2&include=user&space_guids=spaceGUID1", ""},
	}
	setupMultiple(mocks, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	roles, users, err := client.Roles.ListSpaceRolesByGUID("spaceGUID1")
	require.NoError(t, err)
	require.Len(t, roles, 3)
	require.Len(t, users, 3)

	require.Equal(t, "user1", users[0].Username)
	require.Equal(t, "user2", users[1].Username)
	require.Equal(t, "user3", users[2].Username)
}

func TestListSpaceRolesByGUIDAndType(t *testing.T) {
	mocks := []MockRoute{
		{"GET", "/v3/roles", []string{listSpaceRolesBySpaceGUIDAndTypePayload}, "", http.StatusOK, "include=user&space_guids=spaceGUID1&types=space_supporter", ""},
		{"GET", "/v3/rolespage2", []string{listSpaceRolesBySpaceGuidAndTypePayloadPage2}, "", http.StatusOK, "page=2&per_page=2&include=user&space_guids=spaceGUID1&types=space_supporter", ""},
	}
	setupMultiple(mocks, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	users, err := client.Roles.ListSpaceRolesByGUIDAndType("spaceGUID1", "space_supporter")
	require.NoError(t, err)
	require.Len(t, users, 3)

	require.Equal(t, "user1", users[0].Username)
	require.Equal(t, "user2", users[1].Username)
	require.Equal(t, "user3", users[2].Username)
}

func TestListOrgRolesByGUID(t *testing.T) {
	mocks := []MockRoute{
		{"GET", "/v3/roles", []string{listOrganizationRolesByOrganizationGUIDPayload}, "", http.StatusOK, "include=user&organization_guids=orgGUID1", ""},
		{"GET", "/v3/rolespage2", []string{listOrganizationRolesByOrganizationGuidPayloadPage2}, "", http.StatusOK, "page=2&per_page=2&include=user&organization_guids=orgGUID1", ""},
	}
	setupMultiple(mocks, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	roles, users, err := client.Roles.ListOrganizationRolesByGUID("orgGUID1")
	require.NoError(t, err)
	require.Len(t, roles, 3)
	require.Len(t, users, 3)

	require.Equal(t, "user1", users[0].Username)
	require.Equal(t, "user2", users[1].Username)
	require.Equal(t, "user3", users[2].Username)
}

func TestListOrgRolesByGUIDAndType(t *testing.T) {
	mocks := []MockRoute{
		{"GET", "/v3/roles", []string{listOrganizationRolesByOrganizationGUIDAndTypePayload}, "", http.StatusOK, "include=user&organization_guids=orgGUID1&types=organization_auditor", ""},
		{"GET", "/v3/rolespage2", []string{listOrganizationRolesByOrganizationGuidAndTypePayloadPage2}, "", http.StatusOK, "page=2&per_page=2&include=user&organization_guids=orgGUID1&types=organization_auditor", ""},
	}
	setupMultiple(mocks, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	users, err := client.Roles.ListOrganizationRolesByGUIDAndType("orgGUID1", "organization_auditor")
	require.NoError(t, err)
	require.Len(t, users, 3)

	require.Equal(t, "user1", users[0].Username)
	require.Equal(t, "user2", users[1].Username)
	require.Equal(t, "user3", users[2].Username)
}

func TestCreateSpaceRoles(t *testing.T) {
	setup(MockRoute{"POST", "/v3/roles", []string{createSpaceRolePayload}, "", http.StatusCreated, "", ""}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	role, err := client.Roles.CreateSpaceRole(
		"b40a40c8-58b7-49a0-b47d-9d6fe5d72905",
		"c4958204-6b65-43ea-832b-e4c57aea6641",
		"space_supporter",
	)
	require.NoError(t, err)
	require.Equal(t, "space_supporter", role.Type)
	require.Equal(t, "b9f59ab2-2b09-438e-bebb-30e8704ffb89", role.GUID)
	require.Equal(t, "https://api.example.org/v3/roles/b9f59ab2-2b09-438e-bebb-30e8704ffb89", role.Links["self"].Href)
	require.Equal(t, "https://api.example.org/v3/users/c4958204-6b65-43ea-832b-e4c57aea6641", role.Links["user"].Href)
	require.Equal(t, "https://api.example.org/v3/spaces/b40a40c8-58b7-49a0-b47d-9d6fe5d72905", role.Links["space"].Href)
}

func TestCreateOrgRoles(t *testing.T) {
	setup(MockRoute{"POST", "/v3/roles", []string{createOrganizationRolePayload}, "", http.StatusCreated, "", ""}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	role, err := client.Roles.CreateOrganizationRole(
		"fa8a8346-0d92-4729-870c-77ee1934f973",
		"ac2e02c9-2c5c-4712-a620-a68449d263c3",
		"organization_user",
	)
	require.NoError(t, err)
	require.Equal(t, "organization_user", role.Type)
	require.Equal(t, "21cbfaeb-bff7-4cfd-a7a9-6c13ec76f246", role.GUID)
	require.Equal(t, "https://api.example.org/v3/roles/21cbfaeb-bff7-4cfd-a7a9-6c13ec76f246", role.Links["self"].Href)
	require.Equal(t, "https://api.example.org/v3/users/ac2e02c9-2c5c-4712-a620-a68449d263c3", role.Links["user"].Href)
	require.Equal(t, "https://api.example.org/v3/organizations/fa8a8346-0d92-4729-870c-77ee1934f973", role.Links["organization"].Href)
}

func TestDeleteRole(t *testing.T) {
	setup(MockRoute{"DELETE", "/v3/roles/1cb006ee-fb05-47e1-b541-c34179ddc446", []string{""}, "", http.StatusAccepted, "", ""}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	err = client.Roles.Delete("1cb006ee-fb05-47e1-b541-c34179ddc446")
	require.NoError(t, err)
}
