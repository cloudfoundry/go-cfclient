package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"net/http"
	"testing"
)

func TestRoles(t *testing.T) {
	g := test.NewObjectJSONGenerator(15)
	role := g.Role()
	role2 := g.Role()
	role3 := g.Role()
	role4 := g.Role()
	org := g.Organization()
	org2 := g.Organization()
	space := g.Space()
	space2 := g.Space()
	space3 := g.Space()
	user := g.User()
	user2 := g.User()
	user3 := g.User()

	tests := []RouteTest{
		{
			Description: "Create org role",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/roles",
				Output:   []string{role},
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
				return c.Roles.CreateOrganizationRole("ea77cd9e-a072-41e8-9d0b-b2e9180c50bf",
					"0c03442d-c5ae-4661-a929-68f0eeb9ed9a", resource.OrganizationRoleAuditor)
			},
		},
		{
			Description: "Create space role",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/roles",
				Output:   []string{role},
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
				return c.Roles.CreateSpaceRole("c0c8988d-2f97-4768-832a-677557f18174",
					"0c03442d-c5ae-4661-a929-68f0eeb9ed9a", resource.SpaceRoleDeveloper)
			},
		},
		{
			Description: "Get role",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles/211cc662-f86d-4559-a85d-fbfb010c480c",
				Output:   []string{role},
				Status:   http.StatusOK},
			Expected: role,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Roles.Get("211cc662-f86d-4559-a85d-fbfb010c480c")
			},
		},
		{
			Description: "Get role with orgs",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles/211cc662-f86d-4559-a85d-fbfb010c480c",
				Output: g.ResourceWithInclude(test.ResourceResult{
					Resource:      role,
					Organizations: []string{org, org2},
				}),
				Status: http.StatusOK},
			Expected:  role,
			Expected2: g.Array(org, org2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.GetIncludeOrgs("211cc662-f86d-4559-a85d-fbfb010c480c")
			},
		},
		{
			Description: "Get role with spaces",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles/211cc662-f86d-4559-a85d-fbfb010c480c",
				Output: g.ResourceWithInclude(test.ResourceResult{
					Resource: role,
					Spaces:   []string{space, space2},
				}),
				Status: http.StatusOK},
			Expected:  role,
			Expected2: g.Array(space, space2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.GetIncludeSpaces("211cc662-f86d-4559-a85d-fbfb010c480c")
			},
		},
		{
			Description: "Get role with users",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles/211cc662-f86d-4559-a85d-fbfb010c480c",
				Output: g.ResourceWithInclude(test.ResourceResult{
					Resource: role,
					Users:    []string{user, user2, user3},
				}),
				Status: http.StatusOK},
			Expected:  role,
			Expected2: g.Array(user, user2, user3),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.GetIncludeUsers("211cc662-f86d-4559-a85d-fbfb010c480c")
			},
		},
		{
			Description: "Delete role",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/roles/211cc662-f86d-4559-a85d-fbfb010c480c",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Roles.Delete("211cc662-f86d-4559-a85d-fbfb010c480c")
			},
		},
		{
			Description: "List all roles",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles",
				Output:   g.Paged([]string{role, role2}, []string{role3, role4}),
				Status:   http.StatusOK},
			Expected: g.Array(role, role2, role3, role4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Roles.ListAll(nil)
			},
		},
		{
			Description: "List all roles include orgs",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources:     []string{role, role2},
						Organizations: []string{org, org2},
					},
					test.PagedResult{
						Resources:     []string{role3, role4},
						Organizations: []string{},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(role, role2, role3, role4),
			Expected2: g.Array(org, org2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.ListIncludeOrgsAll(nil)
			},
		},
		{
			Description: "List all roles include spaces",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources: []string{role, role2},
						Spaces:    []string{space, space2},
					},
					test.PagedResult{
						Resources: []string{role3, role4},
						Spaces:    []string{space3},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(role, role2, role3, role4),
			Expected2: g.Array(space, space2, space3),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.ListIncludeSpacesAll(nil)
			},
		},
		{
			Description: "List all roles include users",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources: []string{role, role2},
						Users:     []string{user, user2},
					},
					test.PagedResult{
						Resources: []string{role3, role4},
						Users:     []string{user3},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(role, role2, role3, role4),
			Expected2: g.Array(user, user2, user3),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Roles.ListIncludeUsersAll(nil)
			},
		},
	}
	executeTests(tests, t)
}
