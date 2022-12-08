package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestSecurityGroups(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	sg := g.SecurityGroup().JSON
	sg2 := g.SecurityGroup().JSON

	tests := []RouteTest{
		{
			Description: "Bind running security group",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26/relationships/running_spaces",
				Output: []string{`{
					"data": [
						{ "guid": "4ec12cde-e755-4220-9964-65c44c6362b1" },
						{ "guid": "ef498123-7641-44f2-8591-e737c2f96207" }
					],
					"links": {
						"self": {
							"href": "https://api.example.org/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26/relationships/running_spaces"
						}
					}
				}`},
				Status:   http.StatusOK,
				PostForm: `{ "data": [{ "guid": "4ec12cde-e755-4220-9964-65c44c6362b1" }, { "guid": "ef498123-7641-44f2-8591-e737c2f96207" }] }`,
			},
			Expected: `["4ec12cde-e755-4220-9964-65c44c6362b1", "ef498123-7641-44f2-8591-e737c2f96207"]`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SecurityGroups.BindRunningSecurityGroup(context.Background(), "12e9eabb-5139-4377-a5c3-64e3cd1b6e26", []string{
					"4ec12cde-e755-4220-9964-65c44c6362b1", "ef498123-7641-44f2-8591-e737c2f96207",
				})
			},
		},
		{
			Description: "Bind staging security group",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26/relationships/staging_spaces",
				Output: []string{`{
					"data": [
						{ "guid": "4ec12cde-e755-4220-9964-65c44c6362b1" },
						{ "guid": "ef498123-7641-44f2-8591-e737c2f96207" }
					],
					"links": {
						"self": {
							"href": "https://api.example.org/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26/relationships/staging_spaces"
						}
					}
				}`},
				Status:   http.StatusOK,
				PostForm: `{ "data": [{ "guid": "4ec12cde-e755-4220-9964-65c44c6362b1" }, { "guid": "ef498123-7641-44f2-8591-e737c2f96207" }] }`,
			},
			Expected: `["4ec12cde-e755-4220-9964-65c44c6362b1", "ef498123-7641-44f2-8591-e737c2f96207"]`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SecurityGroups.BindStagingSecurityGroup(context.Background(), "12e9eabb-5139-4377-a5c3-64e3cd1b6e26", []string{
					"4ec12cde-e755-4220-9964-65c44c6362b1", "ef498123-7641-44f2-8591-e737c2f96207",
				})
			},
		},
		{
			Description: "Create security group",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/security_groups",
				Output:   g.Single(sg),
				Status:   http.StatusCreated,
				PostForm: `{
				  "name": "my-group0",
				  "rules": [
					{
					  "protocol": "tcp",
					  "destination": "10.10.10.0/24",
					  "ports": "443,80,8080",
					  "log": false
					},
					{
					  "protocol": "icmp",
					  "destination": "10.10.10.0/24",
					  "type": 8,
					  "code": 0,
					  "description": "Allow ping requests to private services"
					}
				  ]
				}`,
			},
			Expected: sg,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.SecurityGroupCreate{
					Name: "my-group0",
					Rules: []*resource.SecurityGroupRule{
						resource.NewSecurityGroupRuleTCP("10.10.10.0/24", false).
							WithPorts("443,80,8080"),
						resource.NewSecurityGroupRuleICMP("10.10.10.0/24", 8, 0).
							WithDescription("Allow ping requests to private services"),
					},
				}
				return c.SecurityGroups.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete security group",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SecurityGroups.Delete(context.Background(), "12e9eabb-5139-4377-a5c3-64e3cd1b6e26")
			},
		},
		{
			Description: "Get security group",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26",
				Output:   g.Single(sg),
				Status:   http.StatusOK},
			Expected: sg,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SecurityGroups.Get(context.Background(), "12e9eabb-5139-4377-a5c3-64e3cd1b6e26")
			},
		},
		{
			Description: "List all security groups",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/security_groups",
				Output:   g.Paged([]string{sg}, []string{sg2}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(sg, sg2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SecurityGroups.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List all running security groups for space",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/spaces/4ec12cde-e755-4220-9964-65c44c6362b1/running_security_groups",
				Output:   g.Paged([]string{sg}, []string{sg2}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(sg, sg2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SecurityGroups.ListRunningForSpaceAll(context.Background(), "4ec12cde-e755-4220-9964-65c44c6362b1", nil)
			},
		},
		{
			Description: "List all staging security groups for space",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/spaces/4ec12cde-e755-4220-9964-65c44c6362b1/staging_security_groups",
				Output:   g.Paged([]string{sg}, []string{sg2}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(sg, sg2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SecurityGroups.ListStagingForSpaceAll(context.Background(), "4ec12cde-e755-4220-9964-65c44c6362b1", nil)
			},
		},
		{
			Description: "Unbind running security group",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26/relationships/running_spaces/4ec12cde-e755-4220-9964-65c44c6362b1",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.SecurityGroups.UnBindRunningSecurityGroup(context.Background(), "12e9eabb-5139-4377-a5c3-64e3cd1b6e26", "4ec12cde-e755-4220-9964-65c44c6362b1")
				return nil, err
			},
		},
		{
			Description: "Unbind staging security group",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26/relationships/staging_spaces/4ec12cde-e755-4220-9964-65c44c6362b1",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.SecurityGroups.UnBindStagingSecurityGroup(context.Background(), "12e9eabb-5139-4377-a5c3-64e3cd1b6e26", "4ec12cde-e755-4220-9964-65c44c6362b1")
				return nil, err
			},
		},
		{
			Description: "Update security group",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26",
				Output:   g.Single(sg),
				Status:   http.StatusOK,
				PostForm: `{
				  "name": "my-group0",
				  "globally_enabled": {
					"running": true,
					"staging": false
				  },
				  "rules": [
					{
					  "protocol": "tcp",
					  "destination": "10.10.10.0/24",
					  "ports": "443,80,8080",
					  "log": false
					},
					{
					  "protocol": "icmp",
					  "destination": "10.10.10.0/24",
					  "type": 8,
					  "code": 0,
					  "description": "Allow ping requests to private services"
					}
				  ]
				}`,
			},
			Expected: sg,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.SecurityGroupUpdate{
					Name: "my-group0",
					GloballyEnabled: &resource.SecurityGroupGloballyEnabled{
						Running: true,
					},
					Rules: []*resource.SecurityGroupRule{
						resource.NewSecurityGroupRuleTCP("10.10.10.0/24", false).
							WithPorts("443,80,8080"),
						resource.NewSecurityGroupRuleICMP("10.10.10.0/24", 8, 0).
							WithDescription("Allow ping requests to private services"),
					},
				}
				return c.SecurityGroups.Update(context.Background(), "12e9eabb-5139-4377-a5c3-64e3cd1b6e26", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
