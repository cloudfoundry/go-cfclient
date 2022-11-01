package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"net/http"
	"testing"
)

func TestRoles(t *testing.T) {
	g := test.NewObjectJSONGenerator(15)
	role := g.Role()
	role2 := g.Role()
	role3 := g.Role()
	role4 := g.Role()

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
			Description: "List paged roles",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/roles",
				Output:   g.Paged([]string{role, role2}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(role, role2),
			Action: func(c *Client, t *testing.T) (any, error) {
				roles, _, err := c.Roles.List(nil)
				return roles, err
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
	}
	executeTests(tests, t)
}
