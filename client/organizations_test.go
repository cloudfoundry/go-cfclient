package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

func TestCreateOrganization(t *testing.T) {
	expectedBody := `{"name":"my-org"}`
	setup(MockRoute{"POST", "/v3/organizations", []string{createOrganizationPayload}, "", http.StatusCreated, "", &expectedBody}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	organization, err := client.CreateOrganization(resource.CreateOrganizationRequest{
		Name: "my-org",
	})
	require.NoError(t, err)
	require.NotNil(t, organization)

	require.Equal(t, "org-guid", organization.GUID)
	require.Equal(t, "quota-guid", organization.Relationships["quota"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/organizations/org-guid/domains", organization.Links["domains"].Href)
	require.Len(t, organization.Metadata.Annotations, 0)
	require.Contains(t, organization.Metadata.Labels, "ORG_KEY")
	require.Equal(t, "org_value", organization.Metadata.Labels["ORG_KEY"])
}

func TestGetOrganization(t *testing.T) {
	setup(MockRoute{"GET", "/v3/organizations/org-guid", []string{getOrganizationPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	organization, err := client.GetOrganizationByGUID("org-guid")
	require.NoError(t, err)
	require.NotNil(t, organization)

	require.Equal(t, "org-guid", organization.GUID)
	require.Equal(t, "quota-guid", organization.Relationships["quota"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/organizations/org-guid/domains", organization.Links["domains"].Href)
	require.Len(t, organization.Metadata.Annotations, 0)
	require.Contains(t, organization.Metadata.Labels, "ORG_KEY")
	require.Equal(t, "org_value", organization.Metadata.Labels["ORG_KEY"])
}

func TestDeleteOrganization(t *testing.T) {
	setup(MockRoute{"DELETE", "/v3/organizations/org-guid", []string{""}, "", http.StatusAccepted, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	err = client.DeleteOrganization("org-guid")
	require.NoError(t, err)
}

func TestUpdateOrganization(t *testing.T) {
	setup(MockRoute{"PATCH", "/v3/organizations/org-guid", []string{updateOrganizationPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	organization, err := client.UpdateOrganization("org-guid", resource.UpdateOrganizationRequest{
		Name: "my-org",
	})
	require.NoError(t, err)
	require.NotNil(t, organization)

	require.Equal(t, "my-org", organization.Name)
	require.Equal(t, "org-guid", organization.GUID)
	require.Equal(t, "", organization.Relationships["quota"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/organizations/org-guid/domains", organization.Links["domains"].Href)
	require.Len(t, organization.Metadata.Annotations, 0)
	require.Contains(t, organization.Metadata.Labels, "ORG_KEY")
	require.Equal(t, "org_value", organization.Metadata.Labels["ORG_KEY"])
}

func TestListOrganizationsByQuery(t *testing.T) {
	setup(MockRoute{"GET", "/v3/organizations", []string{listOrganizationsPayload, listOrganizationsPayloadPage2}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	organizations, err := client.ListOrganizationsByQuery(nil)
	require.NoError(t, err)
	require.Len(t, organizations, 2)

	require.Equal(t, "my-org-1", organizations[0].Name)
	require.Equal(t, "my-org-2", organizations[1].Name)

	require.Equal(t, "quota-guid", organizations[0].Relationships["quota"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/organizations/org-guid/domains", organizations[0].Links["domains"].Href)
	require.Equal(t, "", organizations[1].Relationships["quota"].Data.GUID)
	require.Equal(t, "https://api.example.org/v3/organizations/org-guid-2/domains", organizations[1].Links["domains"].Href)
}
