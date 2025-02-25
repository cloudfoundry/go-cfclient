package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/cloudfoundry/go-cfclient/v3/testutil"
)

func TestEnvVarGroups(t *testing.T) {
	g := testutil.NewObjectJSONGenerator()
	envVarGroup := g.EnvVarGroup().JSON

	tests := []RouteTest{
		{
			Description: "Get running env var group",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/environment_variable_groups/running",
				Output:   g.Single(envVarGroup),
				Status:   http.StatusOK},
			Expected: envVarGroup,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.EnvVarGroups.GetRunning(context.Background())
			},
		},
		{
			Description: "Update buildpack",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/environment_variable_groups/staging",
				Output:   g.Single(envVarGroup),
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
				return c.EnvVarGroups.UpdateStaging(context.Background(), r)
			},
		},
	}
	ExecuteTests(tests, t)
}
