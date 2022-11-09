package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"net/http"
	"testing"
)

func TestEnvVarGroups(t *testing.T) {
	g := test.NewObjectJSONGenerator(852)
	envVarGroup := g.EnvVarGroup()

	tests := []RouteTest{
		{
			Description: "Get running env var group",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/environment_variable_groups/running",
				Output:   []string{envVarGroup},
				Status:   http.StatusOK},
			Expected: envVarGroup,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.EnvVarGroups.GetRunning()
			},
		},
		{
			Description: "Update buildpack",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/environment_variable_groups/staging",
				Output:   []string{envVarGroup},
				Status:   http.StatusOK,
				PostForm: `{ "var": { "DEBUG": "false" }}`,
			},
			Expected: envVarGroup,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.EnvVarGroupUpdate{
					Var: map[string]string{
						"DEBUG": "false",
					},
				}
				return c.EnvVarGroups.UpdateStaging(r)
			},
		},
	}
	executeTests(tests, t)
}
