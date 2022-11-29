package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestRoles(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(15)
	role := g.Role().JSON
	role2 := g.Role().JSON
	role3 := g.Role().JSON
	role4 := g.Role().JSON
	org := g.Organization().JSON
	org2 := g.Organization().JSON
	space := g.Space().JSON
	space2 := g.Space().JSON
	space3 := g.Space().JSON
	user := g.User().JSON
	user2 := g.User().JSON
	user3 := g.User().JSON

	tests := []RouteTest{
		{
			Description: "Create organization role",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/roles",
				Output:   g.Single(role),
				Status:   http.StatusCreated,
				PostForm: `{
				  "type": "organization_auditor",
				  "relationships": {
					"user": {
					  "data": {
						"guid": "0c03442d-c5ae-4661-a929-68f0eeb9ed9a"
					  }
					},
					"organization": {
					  "data": {
						"guid": "ea77cd9e-a072-41e8-9d0b-b2e9180c50bf"
					  }
					}
				  }
				}`,
			},
			Expected: role,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Roles.CreateOrganizationRole(context.Background(), "ea77cd9e-a072-41e8-9d0b-b2e9180c50bf",
					"0c03442d-c5ae-4661-a929-68f0eeb9ed9a", resource.OrganizationRoleAuditor)
			},
		},
		{
			Description: "Create space role",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/roles",
				Output:   g.Single(role),
				Status:   http.StatusCreated,
				PostForm: `{
				  "type": "space_developer",
				  "relationships": {
					"user": {
					  "data": {
						"guid": "0c03442d-c5ae-4661-a929-68f0eeb9ed9a"
					  }
					},
					"space": {
					  "data": {
						"guid": "c0c8988d-2f97-4768-832a-677557f18174"
					  }
					}
				  }
				}`,
			},
			Expected: role,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Roles.CreateSpaceRole(context.Background(), "c0c8988d-2f97-4768-832a-677557f18174",
					"0c03442d-c5ae-4661-a929-68f0eeb9ed9a", resource.SpaceRoleDeveloper)
			},
		},
		{
			Description: "Delete role",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/roles/211cc662-f86d-4559-a85d-fbfb010c480c",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Roles.Delete(context.Background(), "211cc662-f86d-4559-a85d-fbfb010c480c")
			},
		},
		{
			Description: "Get role",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles/211cc662-f86d-4559-a85d-fbfb010c480c",
				Output:   g.Single(role),
				Status:   http.StatusOK},
			Expected: role,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Roles.Get(context.Background(), "211cc662-f86d-4559-a85d-fbfb010c480c")
			},
		},
		{
			Description: "Get role with organizations",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles/211cc662-f86d-4559-a85d-fbfb010c480c",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource:      role,
					Organizations: []string{org, org2},
				}),
				Status: http.StatusOK},
			Expected:  role,
			Expected2: g.Array(org, org2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.GetIncludeOrganizations(context.Background(), "211cc662-f86d-4559-a85d-fbfb010c480c")
			},
		},
		{
			Description: "Get role with spaces",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles/211cc662-f86d-4559-a85d-fbfb010c480c",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource: role,
					Spaces:   []string{space, space2},
				}),
				Status: http.StatusOK},
			Expected:  role,
			Expected2: g.Array(space, space2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.GetIncludeSpaces(context.Background(), "211cc662-f86d-4559-a85d-fbfb010c480c")
			},
		},
		{
			Description: "Get role with users",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles/211cc662-f86d-4559-a85d-fbfb010c480c",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource: role,
					Users:    []string{user, user2, user3},
				}),
				Status: http.StatusOK},
			Expected:  role,
			Expected2: g.Array(user, user2, user3),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.GetIncludeUsers(context.Background(), "211cc662-f86d-4559-a85d-fbfb010c480c")
			},
		},
		{
			Description: "List all roles",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles",
				Output:   g.Paged([]string{role, role2}, []string{role3, role4}),
				Status:   http.StatusOK},
			Expected: g.Array(role, role2, role3, role4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Roles.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List all roles include organizations",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources:     []string{role, role2},
						Organizations: []string{org, org2},
					},
					testutil.PagedResult{
						Resources:     []string{role3, role4},
						Organizations: []string{},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(role, role2, role3, role4),
			Expected2: g.Array(org, org2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.ListIncludeOrganizationsAll(context.Background(), nil)
			},
		},
		{
			Description: "List all roles include spaces",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{role, role2},
						Spaces:    []string{space, space2},
					},
					testutil.PagedResult{
						Resources: []string{role3, role4},
						Spaces:    []string{space3},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(role, role2, role3, role4),
			Expected2: g.Array(space, space2, space3),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.ListIncludeSpacesAll(context.Background(), nil)
			},
		},
		{
			Description: "List all roles include users",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{role, role2},
						Users:     []string{user, user2},
					},
					testutil.PagedResult{
						Resources: []string{role3, role4},
						Users:     []string{user3},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(role, role2, role3, role4),
			Expected2: g.Array(user, user2, user3),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.ListIncludeUsersAll(context.Background(), nil)
			},
		},
	}
	ExecuteTests(tests, t)
}
