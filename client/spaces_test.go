package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

func TestCreateSpace(t *testing.T) {
	expectedBody := `{"name":"my-space","relationships":{"organization":{"data":{"guid":"org-guid"}}}}`
	setup(MockRoute{"POST", "/v3/spaces", []string{createSpacePayload}, "", http.StatusCreated, "", &expectedBody}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	space, err := client.Spaces.Create(resource.CreateSpaceRequest{
		Name:    "my-space",
		OrgGUID: "org-guid",
	})
	require.NoError(t, err)
	require.NotNil(t, space)

	require.Equal(t, "space-guid", space.GUID)
	require.Equal(t, "org-guid", space.Relationships["organization"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/organizations/org-guid", space.Links["organization"].Href)
	require.Len(t, space.Metadata.Annotations, 0)
	require.Contains(t, space.Metadata.Labels, "SPACE_KEY")
	require.Equal(t, "space_value", space.Metadata.Labels["SPACE_KEY"])
}

func TestGetSpace(t *testing.T) {
	setup(MockRoute{"GET", "/v3/spaces/space-guid", []string{getSpacePayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	space, err := client.Spaces.Get("space-guid")
	require.NoError(t, err)
	require.NotNil(t, space)

	require.Equal(t, "space-guid", space.GUID)
	require.Equal(t, "org-guid", space.Relationships["organization"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/organizations/org-guid", space.Links["organization"].Href)
	require.Len(t, space.Metadata.Annotations, 0)
	require.Contains(t, space.Metadata.Labels, "SPACE_KEY")
	require.Equal(t, "space_value", space.Metadata.Labels["SPACE_KEY"])
}

func TestDeleteSpace(t *testing.T) {
	setup(MockRoute{"DELETE", "/v3/spaces/space-guid", []string{""}, "", http.StatusAccepted, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	err = client.Spaces.Delete("space-guid")
	require.NoError(t, err)
}

func TestUpdateSpace(t *testing.T) {
	setup(MockRoute{"PATCH", "/v3/spaces/space-guid", []string{updateSpacePayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	space, err := client.Spaces.Update("space-guid", resource.UpdateSpaceRequest{
		Name: "my-space",
	})
	require.NoError(t, err)
	require.NotNil(t, space)

	require.Equal(t, "my-space", space.Name)
	require.Equal(t, "space-guid", space.GUID)
	require.Equal(t, "org-guid", space.Relationships["organization"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/organizations/org-guid", space.Links["organization"].Href)
	require.Len(t, space.Metadata.Annotations, 0)
	require.Contains(t, space.Metadata.Labels, "SPACE_KEY")
	require.Equal(t, "space_value", space.Metadata.Labels["SPACE_KEY"])
}

func TestListSpacesByQuery(t *testing.T) {
	setup(MockRoute{"GET", "/v3/spaces", []string{listSpacesPayload, listSpacesPayloadPage2}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	spaces, err := client.Spaces.ListByQuery(nil)
	require.NoError(t, err)
	require.Len(t, spaces, 2)

	require.Equal(t, "my-space-1", spaces[0].Name)
	require.Equal(t, "my-space-2", spaces[1].Name)

	require.Equal(t, "org-guid", spaces[0].Relationships["organization"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/organizations/org-guid", spaces[0].Links["organization"].Href)
	require.Equal(t, "org-guid", spaces[1].Relationships["organization"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/organizations/org-guid", spaces[1].Links["organization"].Href)
}

func TestListSpaceUsersByQuery(t *testing.T) {
	setup(MockRoute{"GET", "/v3/spaces/space-guid/users", []string{listSpaceUsersPayload, listSpaceUsersPayloadPage2}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	users, err := client.Spaces.ListUsers("space-guid")
	require.NoError(t, err)
	require.Len(t, users, 2)

	require.Equal(t, "some-name-1", users[0].Username)
	require.Equal(t, "some-name-2", users[1].Username)

	require.Equal(t, "some-name-1", users[0].PresentationName)
	require.Equal(t, "uaa", users[0].Origin)
	require.Equal(t, "https://api.example.org/v3/users/10a93b89-3f89-4f05-7238-8a2b123c79l9", users[0].Links["self"].Href)
	require.Equal(t, "some-name-2", users[1].PresentationName)
	require.Equal(t, "ldap", users[1].Origin)
	require.Equal(t, "https://api.example.org/v3/users/9da93b89-3f89-4f05-7238-8a2b123c79l9", users[1].Links["self"].Href)
}
