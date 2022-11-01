package client

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestSecurityGroups(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	sg := g.SecurityGroup()
	sg2 := g.SecurityGroup()

	tests := []RouteTest{
		{
			Description: "Create security group",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/security_groups",
				Output:   []string{sg},
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
				return c.SecurityGroups.Create(r)
			},
		},
		{
			Description: "Delete security group",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.SecurityGroups.Delete("12e9eabb-5139-4377-a5c3-64e3cd1b6e26")
			},
		},
		{
			Description: "Get security group",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26",
				Output:   []string{sg},
				Status:   http.StatusOK},
			Expected: sg,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SecurityGroups.Get("12e9eabb-5139-4377-a5c3-64e3cd1b6e26")
			},
		},
		{
			Description: "List all security groups",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/security_groups",
				Output:   g.Paged([]string{sg}, []string{sg2}),
				Status:   http.StatusOK},
			Expected: g.Array(sg, sg2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SecurityGroups.ListAll(nil)
			},
		},
		{
			Description: "Update security group",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/security_groups/12e9eabb-5139-4377-a5c3-64e3cd1b6e26",
				Output:   []string{sg},
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
				return c.SecurityGroups.Update("12e9eabb-5139-4377-a5c3-64e3cd1b6e26", r)
			},
		},
	}
	for _, tt := range tests {
		func() {
			setup(tt.Route, t)
			defer teardown()
			details := fmt.Sprintf("%s %s", tt.Route.Method, tt.Route.Endpoint)
			if tt.Description != "" {
				details = tt.Description + ": " + details
			}

			c, _ := NewTokenConfig(server.URL, "foobar")
			cl, err := New(c)
			require.NoError(t, err, details)

			obj, err := tt.Action(cl, t)
			require.NoError(t, err, details)
			if tt.Expected != "" {
				actual, err := json.Marshal(obj)
				require.NoError(t, err, details)
				require.JSONEq(t, tt.Expected, string(actual), details)
			}
		}()
	}
}
