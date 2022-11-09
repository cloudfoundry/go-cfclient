package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"net/http"
	"testing"
)

func TestApps(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	app1 := g.Application()
	app2 := g.Application()
	app3 := g.Application()
	app4 := g.Application()
	space1 := g.Space()
	space2 := g.Space()
	org := g.Organization()
	appEnvironment := g.AppEnvironment()
	appEnvVar := g.AppEnvVar()
	appSSH := g.AppSSH()
	appPermission := g.AppPermission()

	tests := []RouteTest{
		{
			Description: "Create app",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/apps",
				Output:   []string{app1},
				Status:   http.StatusCreated,
				PostForm: `{"environment_variables":{"FOO":"BAR"},"name":"my-app","relationships":{"space":{"data":{"guid":"space-guid"}}}}`,
			},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewAppCreate("my-app", "space-guid")
				r.EnvironmentVariables = map[string]string{"FOO": "BAR"}
				return c.Applications.Create(r)
			},
		},
		{
			Description: "Get app",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446",
				Output:   []string{app1},
				Status:   http.StatusOK},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Get("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Get app environment",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/env",
				Output:   []string{appEnvironment},
				Status:   http.StatusOK,
			},
			Expected: appEnvironment,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.GetEnvironment("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Update app environment variables",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables",
				Output:   []string{g.AppUpdateEnvVars()},
				Status:   http.StatusOK,
			},
			Expected: `{ "RAILS_ENV": "production", "DEBUG": "false" }`,
			Action: func(c *Client, t *testing.T) (any, error) {
				falseVar := "false"
				return c.Applications.SetEnvironmentVariables("1cb006ee-fb05-47e1-b541-c34179ddc446",
					map[string]*string{
						"DEBUG": &falseVar,
						"USER":  nil,
					},
				)
			},
		},
		{
			Description: "Get app environment variables",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables",
				Output:   []string{appEnvVar},
				Status:   http.StatusOK,
			},
			Expected: `{ "RAILS_ENV": "production" }`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.GetEnvironmentVariables("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Get SSH enabled for app",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/ssh_enabled",
				Output:   []string{appSSH},
				Status:   http.StatusOK,
			},
			Expected: appSSH,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.SSHEnabled("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Get app permissions",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/permissions",
				Output:   []string{appPermission},
				Status:   http.StatusOK,
			},
			Expected: appPermission,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Permissions("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Start app",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/start",
				Output:   []string{app1},
				Status:   http.StatusOK,
			},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Start("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Stop app",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/stop",
				Output:   []string{app1},
				Status:   http.StatusOK,
			},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Stop("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Restart app",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/restart",
				Output:   []string{app1},
				Status:   http.StatusOK,
			},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Restart("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Delete app",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Applications.Delete("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Update app",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446",
				Output:   []string{app1},
				Status:   http.StatusOK,
				PostForm: `{ "name": "new_name", "lifecycle": { "type": "buildpack", "data": { "buildpacks": ["java_offline"] }}}`,
			},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.AppUpdate{
					Name: "new_name",
					Lifecycle: &resource.Lifecycle{
						Type: "buildpack",
						BuildpackData: resource.BuildpackLifecycle{
							Buildpacks: []string{"java_offline"},
						},
					},
				}
				return c.Applications.Update("1cb006ee-fb05-47e1-b541-c34179ddc446", r)
			},
		},
		{
			Description: "List all apps",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps",
				Output:   g.Paged([]string{app1, app2}, []string{app3, app4}),
				Status:   http.StatusOK},
			Expected: g.Array(app1, app2, app3, app4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.ListAll(nil)
			},
		},
		{
			Description: "List all apps include spaces",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources: []string{app1, app2},
						Spaces:    []string{space1},
					},
					test.PagedResult{
						Resources: []string{app3, app4},
						Spaces:    []string{space2},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(app1, app2, app3, app4),
			Expected2: g.Array(space1, space2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Applications.ListIncludeSpacesAll(nil)
			},
		},
		{
			Description: "List all apps include spaces and orgs",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources:     []string{app1, app2},
						Spaces:        []string{space1},
						Organizations: []string{org},
					},
					test.PagedResult{
						Resources: []string{app3, app4},
						Spaces:    []string{space2},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(app1, app2, app3, app4),
			Expected2: g.Array(space1, space2),
			Expected3: g.Array(org),
			Action3: func(c *Client, t *testing.T) (any, any, any, error) {
				return c.Applications.ListIncludeSpacesAndOrgsAll(nil)
			},
		},
	}
	executeTests(tests, t)
}
