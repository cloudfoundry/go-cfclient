package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestOrganizations(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(15)
	org := g.Organization().JSON
	org2 := g.Organization().JSON
	org3 := g.Organization().JSON
	org4 := g.Organization().JSON
	domain := g.Domain().JSON
	orgUsageSummary := g.OrganizationUsageSummary().JSON
	user := g.User().JSON
	user2 := g.User().JSON

	tests := []RouteTest{
		{
			Description: "Assign default organization iso segment",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5/relationships/default_isolation_segment",
				Output:   []string{`{ "data": { "guid": "443a1ea0-2403-4f0f-8c74-023a320bd1f2" }}`},
				Status:   http.StatusOK,
				PostForm: `{ "data": { "guid": "443a1ea0-2403-4f0f-8c74-023a320bd1f2" }}`,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.Organizations.AssignDefaultIsolationSegment(context.Background(), "3691e277-eb88-4ddc-bec3-0111d9dd4ef5", "443a1ea0-2403-4f0f-8c74-023a320bd1f2")
				return nil, err
			},
		},
		{
			Description: "Create organization",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/organizations",
				Output:   g.Single(org),
				Status:   http.StatusCreated,
				PostForm: `{ "name": "my-organization" }`,
			},
			Expected: org,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewOrganizationCreate("my-organization")
				return c.Organizations.Create(context.Background(), r)
			},
		},
		{
			Description: "Get organization",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5",
				Output:   g.Single(org),
				Status:   http.StatusOK,
			},
			Expected: org,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.Get(context.Background(), "3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
			},
		},
		{
			Description: "Get organization default iso segment",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5/relationships/default_isolation_segment",
				Output:   []string{`{ "data": { "guid": "443a1ea0-2403-4f0f-8c74-023a320bd1f2" }}`},
				Status:   http.StatusOK,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				iso, err := c.Organizations.GetDefaultIsolationSegment(context.Background(), "3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
				require.NoError(t, err)
				require.Equal(t, "443a1ea0-2403-4f0f-8c74-023a320bd1f2", iso)
				return nil, nil
			},
		},
		{
			Description: "Get organization default domain",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5/domains/default",
				Output:   g.Single(domain),
				Status:   http.StatusOK,
			},
			Expected: domain,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.GetDefaultDomain(context.Background(), "3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
			},
		},
		{
			Description: "Get organization usage summary",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5/usage_summary",
				Output:   g.Single(orgUsageSummary),
				Status:   http.StatusOK,
			},
			Expected: orgUsageSummary,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.GetUsageSummary(context.Background(), "3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
			},
		},
		{
			Description: "Delete organization",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.Delete(context.Background(), "3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
			},
		},
		{
			Description: "List all organizations",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations",
				Output:   g.Paged([]string{org, org2}, []string{org3, org4}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(org, org2, org3, org4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List all organizations for iso segment",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/isolation_segments/571de34f-8067-44f0-8bec-4ac17bf8750f/organizations",
				Output:   g.Paged([]string{org, org2}, []string{org3, org4}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(org, org2, org3, org4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.ListForIsolationSegmentAll(context.Background(), "571de34f-8067-44f0-8bec-4ac17bf8750f", nil)
			},
		},
		{
			Description: "List all organization users",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5/users",
				Output:   g.Paged([]string{user}, []string{user2}),
				Status:   http.StatusOK},
			Expected: g.Array(user, user2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.ListUsersAll(context.Background(), "3691e277-eb88-4ddc-bec3-0111d9dd4ef5", nil)
			},
		},
		{
			Description: "Update organization",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5",
				Output:   g.Single(org),
				Status:   http.StatusOK,
				PostForm: `{ "name": "new_name" }`,
			},
			Expected: org,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.OrganizationUpdate{
					Name: "new_name",
				}
				return c.Organizations.Update(context.Background(), "3691e277-eb88-4ddc-bec3-0111d9dd4ef5", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
