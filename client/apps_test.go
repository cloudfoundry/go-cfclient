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

func TestApps(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	app1 := g.Application()
	app2 := g.Application()
	app3 := g.Application()
	app4 := g.Application()
	appEnvVars := g.AppEnvVars()
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
				PostForm: `{"environment_variables":{"FOO":"BAR"},"name":"my-app","relationships":{"space":{"data":{"guid":"space-guid"}}}}`},
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
			Description: "Get app env vars",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/env",
				Output:   []string{appEnvVars},
				Status:   http.StatusOK,
			},
			Expected: appEnvVars,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.GetEnvironment("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Update app env vars",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables",
				Output:   []string{g.AppUpdateEnvVars()},
				Status:   http.StatusOK,
			},
			Expected: `{"var": { "RAILS_ENV": "production", "DEBUG": "false" }}`,
			Action: func(c *Client, t *testing.T) (any, error) {
				falseVar := "false"
				return c.Applications.SetEnvVariables("1cb006ee-fb05-47e1-b541-c34179ddc446",
					resource.EnvVar{Var: map[string]*string{
						"DEBUG": &falseVar,
						"USER":  nil,
					}},
				)
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
			},
			Expected: app1,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.Update("1cb006ee-fb05-47e1-b541-c34179ddc446", &resource.AppUpdate{})
			},
		},
		{
			Description: "List first page of apps",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps",
				Output:   g.Paged("apps", []string{app1}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(app1),
			Action: func(c *Client, t *testing.T) (any, error) {
				apps, _, err := c.Applications.List(NewAppListOptions())
				return apps, err
			},
		},
		{
			Description: "List all apps",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps",
				Output:   g.Paged("apps", []string{app1, app2}, []string{app3, app4}),
				Status:   http.StatusOK},
			Expected: g.Array(app1, app2, app3, app4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Applications.ListAll()
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
