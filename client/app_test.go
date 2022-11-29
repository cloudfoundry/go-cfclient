package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestApps(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	app1 := g.Application().JSON
	app2 := g.Application().JSON
	app3 := g.Application().JSON
	app4 := g.Application().JSON
	space1 := g.Space().JSON
	space2 := g.Space().JSON
	org := g.Organization().JSON
	appEnvironment := g.AppEnvironment().JSON
	appEnvVar := g.AppEnvVar().JSON
	appSSH := g.AppSSH().JSON
	appPermission := g.AppPermission().JSON

	tests := []RouteTest{
		{
			Description: "Create app",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/apps",
				Output:   g.Single(app1),
				Status:   http.StatusCreated,
				PostForm: `{"environment_variables":{"FOO":"BAR"},"name":"my-app","relationships":{"space":{"data":{"guid":"space-guid"}}}}`,
			},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewAppCreate("my-app", "space-guid")
				r.EnvironmentVariables = map[string]string{"FOO": "BAR"}
				return c.Applications.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete app",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Delete(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "first app",
			Route: testutil.MockRoute{
				Method:      "GET",
				Endpoint:    "/v3/apps",
				QueryString: "names=spring-music&page=1&per_page=50",
				Output:      g.Paged([]string{app1, app2}),
				Status:      http.StatusOK},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				opts := NewAppListOptions()
				opts.Names.EqualTo("spring-music")
				return c.Applications.First(context.Background(), opts)
			},
		},
		{
			Description: "first app matches 0 apps",
			Route: testutil.MockRoute{
				Method:      "GET",
				Endpoint:    "/v3/apps",
				QueryString: "names=spring-music&page=1&per_page=50",
				Output:      g.Paged([]string{}),
				Status:      http.StatusOK},
			Action: func(c *Client, t *testing.T) (any, error) {
				opts := NewAppListOptions()
				opts.Names.EqualTo("spring-music")
				app, err := c.Applications.First(context.Background(), opts)
				require.Nil(t, app)
				require.Same(t, ErrNoResultsReturned, err)
				return nil, nil
			},
		},
		{
			Description: "Get app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446",
				Output:   g.Single(app1),
				Status:   http.StatusOK},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Get(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Get app environment",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/env",
				Output:   g.Single(appEnvironment),
				Status:   http.StatusOK,
			},
			Expected: appEnvironment,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.GetEnvironment(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Update app environment variables",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables",
				Output:   []string{g.AppUpdateEnvVars().JSON},
				Status:   http.StatusOK,
			},
			Expected: `{ "RAILS_ENV": "production", "DEBUG": "false" }`,
			Action: func(c *Client, t *testing.T) (any, error) {
				falseVar := "false"
				return c.Applications.SetEnvironmentVariables(context.Background(),
					"1cb006ee-fb05-47e1-b541-c34179ddc446",
					map[string]*string{
						"DEBUG": &falseVar,
						"USER":  nil,
					},
				)
			},
		},
		{
			Description: "Get app environment variables",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables",
				Output:   g.Single(appEnvVar),
				Status:   http.StatusOK,
			},
			Expected: `{ "RAILS_ENV": "production" }`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.GetEnvironmentVariables(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Get SSH enabled for app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/ssh_enabled",
				Output:   g.Single(appSSH),
				Status:   http.StatusOK,
			},
			Expected: appSSH,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.SSHEnabled(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Get app permissions",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/permissions",
				Output:   g.Single(appPermission),
				Status:   http.StatusOK,
			},
			Expected: appPermission,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Permissions(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "List all apps",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps",
				Output:   g.Paged([]string{app1, app2}, []string{app3, app4}),
				Status:   http.StatusOK},
			Expected: g.Array(app1, app2, app3, app4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List all apps include spaces",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{app1, app2},
						Spaces:    []string{space1},
					},
					testutil.PagedResult{
						Resources: []string{app3, app4},
						Spaces:    []string{space2},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(app1, app2, app3, app4),
			Expected2: g.Array(space1, space2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Applications.ListIncludeSpacesAll(context.Background(), nil)
			},
		},
		{
			Description: "List all apps include spaces and organizations",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources:     []string{app1, app2},
						Spaces:        []string{space1},
						Organizations: []string{org},
					},
					testutil.PagedResult{
						Resources: []string{app3, app4},
						Spaces:    []string{space2},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(app1, app2, app3, app4),
			Expected2: g.Array(space1, space2),
			Expected3: g.Array(org),
			Action3: func(c *Client, t *testing.T) (any, any, any, error) {
				return c.Applications.ListIncludeSpacesAndOrganizationsAll(context.Background(), nil)
			},
		},
		{
			Description: "single app",
			Route: testutil.MockRoute{
				Method:      "GET",
				Endpoint:    "/v3/apps",
				QueryString: "names=spring-music-124&page=1&per_page=50",
				Output:      g.Paged([]string{app1}),
				Status:      http.StatusOK},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				opts := NewAppListOptions()
				opts.Names.EqualTo("spring-music-124")
				return c.Applications.Single(context.Background(), opts)
			},
		},
		{
			Description: "single app matches 2+ apps",
			Route: testutil.MockRoute{
				Method:      "GET",
				Endpoint:    "/v3/apps",
				QueryString: "names=spring-music&page=1&per_page=50",
				Output:      g.Paged([]string{app1, app2}),
				Status:      http.StatusOK},
			Action: func(c *Client, t *testing.T) (any, error) {
				opts := NewAppListOptions()
				opts.Names.EqualTo("spring-music")
				app, err := c.Applications.Single(context.Background(), opts)
				require.Nil(t, app)
				require.Same(t, ErrExactlyOneResultNotReturned, err)
				return nil, nil
			},
		},
		{
			Description: "Start app",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/start",
				Output:   g.Single(app1),
				Status:   http.StatusOK,
			},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Start(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Stop app",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/stop",
				Output:   g.Single(app1),
				Status:   http.StatusOK,
			},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Stop(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Restart app",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/restart",
				Output:   g.Single(app1),
				Status:   http.StatusOK,
			},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Restart(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Update app",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446",
				Output:   g.Single(app1),
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
				return c.Applications.Update(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
