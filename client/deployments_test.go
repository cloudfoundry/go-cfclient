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

func TestDeployments(t *testing.T) {
	g := test.NewObjectJSONGenerator(4)
	deployment := g.Deployment()
	deployment2 := g.Deployment()
	deployment3 := g.Deployment()
	deployment4 := g.Deployment()

	tests := []RouteTest{
		{
			Description: "Create deployment with droplet",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/deployments",
				Output:   []string{deployment},
				Status:   http.StatusCreated,
				PostForm: `{"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}, "droplet": {"guid": "c2941033-4575-486d-bf2c-3ae49e8b4ca1"}}`,
			},
			Expected: deployment,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewDeploymentCreate("305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
				r.Droplet = &resource.Relationship{
					GUID: "c2941033-4575-486d-bf2c-3ae49e8b4ca1",
				}
				return c.Deployments.Create(r)
			},
		},
		{
			Description: "Create deployment with revision",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/deployments",
				Output:   []string{deployment},
				Status:   http.StatusCreated,
				PostForm: `{"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}, "revision": {"guid": "d95d8024-8665-4aac-97ea-3c08373e233e"}}`,
			},
			Expected: deployment,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewDeploymentCreate("305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
				r.Revision = &resource.DeploymentRevision{
					GUID: "d95d8024-8665-4aac-97ea-3c08373e233e",
				}
				return c.Deployments.Create(r)
			},
		},
		{
			Description: "Create deployment with revision and droplet",
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewDeploymentCreate("305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
				r.Revision = &resource.DeploymentRevision{
					GUID: "d95d8024-8665-4aac-97ea-3c08373e233e",
				}
				r.Droplet = &resource.Relationship{
					GUID: "c2941033-4575-486d-bf2c-3ae49e8b4ca1",
				}
				_, err := c.Deployments.Create(r)
				require.Error(t, err)
				require.ErrorContains(t, err, "droplet and revision cannot both be set")
				return nil, nil
			},
		},
		{
			Description: "Cancel deployment",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/deployments/2b56dc7b-2a14-49ea-be29-ca182b14a998/actions/cancel",
				Status:   http.StatusOK,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Deployments.Cancel("2b56dc7b-2a14-49ea-be29-ca182b14a998")
			},
		},
		{
			Description: "Get deployment",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/deployments/2b56dc7b-2a14-49ea-be29-ca182b14a998",
				Output:   []string{deployment},
				Status:   http.StatusOK,
			},
			Expected: deployment,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Deployments.Get("2b56dc7b-2a14-49ea-be29-ca182b14a998")
			},
		},
		{
			Description: "List first page of deployments",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/deployments",
				Output:   g.Paged("deployments", []string{deployment}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(deployment),
			Action: func(c *Client, t *testing.T) (any, error) {
				apps, _, err := c.Deployments.List(NewDeploymentListOptions())
				return apps, err
			},
		},
		{
			Description: "List all apps",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/deployments",
				Output:   g.Paged("deployments", []string{deployment, deployment2}, []string{deployment3, deployment4}),
				Status:   http.StatusOK},
			Expected: g.Array(deployment, deployment2, deployment3, deployment4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Deployments.ListAll()
			},
		},
		{
			Description: "Update deployment",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/deployments/2b56dc7b-2a14-49ea-be29-ca182b14a998",
				Output:   []string{deployment},
				PostForm: `{ "metadata": { "labels": { "key": "value" }, "annotations": {"note": "detailed information"}}}`,
				Status:   http.StatusOK,
			},
			Expected: deployment,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.DeploymentUpdate{
					Metadata: &resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.Deployments.Update("2b56dc7b-2a14-49ea-be29-ca182b14a998", r)
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
