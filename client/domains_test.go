package client

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestListDomains(t *testing.T) {
	setup(MockRoute{"GET", "/v3/domains", []string{listDomainsPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	resp, err := client.Domains.ListByQuery(nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Len(t, resp, 1)
	require.Equal(t, "test-domain.com", resp[0].Name)
	require.Equal(t, "3a5d3d89-3f89-4f05-8188-8a2b298c79d5", resp[0].GUID)
	require.Equal(t, false, resp[0].Internal)
	require.Equal(t, "3a3f3d89-3f89-4f05-8188-751b298c79d5", resp[0].Relationships.Organization.Data.GUID)
	require.Equal(t, "404f3d89-3f89-6z72-8188-751b298d88d5", resp[0].Relationships.SharedOrganizations.Data[0].GUID)
	require.Equal(t, "416d3d89-3f89-8h67-2189-123b298d3592", resp[0].Relationships.SharedOrganizations.Data[1].GUID)
	require.Equal(t, "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5", resp[0].Links["self"].Href)
	require.Equal(t, "https://api.example.org/v3/organizations/3a3f3d89-3f89-4f05-8188-751b298c79d5", resp[0].Links["organization"].Href)
	require.Equal(t, "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5/route_reservations", resp[0].Links["route_reservations"].Href)
	require.Equal(t, "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5/relationships/shared_organizations", resp[0].Links["shared_organizations"].Href)
}
