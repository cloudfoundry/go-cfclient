package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestDeployments(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(4)
	deployment := g.Deployment().JSON
	deployment2 := g.Deployment().JSON
	deployment3 := g.Deployment().JSON
	deployment4 := g.Deployment().JSON

	tests := []RouteTest{
		{
			Description: "Create deployment with droplet",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/deployments",
				Output:   g.Single(deployment),
				Status:   http.StatusCreated,
				PostForm: `{"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}, "droplet": {"guid": "c2941033-4575-486d-bf2c-3ae49e8b4ca1"}}`,
			},
			Expected: deployment,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewDeploymentCreate("305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
				r.Droplet = &resource.Relationship{
					GUID: "c2941033-4575-486d-bf2c-3ae49e8b4ca1",
				}
				return c.Deployments.Create(context.Background(), r)
			},
		},
		{
			Description: "Create deployment with revision",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/deployments",
				Output:   g.Single(deployment),
				Status:   http.StatusCreated,
				PostForm: `{"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}, "revision": {"guid": "d95d8024-8665-4aac-97ea-3c08373e233e"}}`,
			},
			Expected: deployment,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewDeploymentCreate("305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
				r.Revision = &resource.DeploymentRevision{
					GUID: "d95d8024-8665-4aac-97ea-3c08373e233e",
				}
				return c.Deployments.Create(context.Background(), r)
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
				_, err := c.Deployments.Create(context.Background(), r)
				require.Error(t, err)
				require.ErrorContains(t, err, "droplet and revision cannot both be set")
				return nil, nil
			},
		},
		{
			Description: "Cancel deployment",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/deployments/2b56dc7b-2a14-49ea-be29-ca182b14a998/actions/cancel",
				Status:   http.StatusOK,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Deployments.Cancel(context.Background(), "2b56dc7b-2a14-49ea-be29-ca182b14a998")
			},
		},
		{
			Description: "Get deployment",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/deployments/2b56dc7b-2a14-49ea-be29-ca182b14a998",
				Output:   g.Single(deployment),
				Status:   http.StatusOK,
			},
			Expected: deployment,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Deployments.Get(context.Background(), "2b56dc7b-2a14-49ea-be29-ca182b14a998")
			},
		},
		{
			Description: "List first page of deployments",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/deployments",
				Output:   g.SinglePaged(deployment),
				Status:   http.StatusOK,
			},
			Expected: g.Array(deployment),
			Action: func(c *Client, t *testing.T) (any, error) {
				apps, _, err := c.Deployments.List(context.Background(), nil)
				return apps, err
			},
		},
		{
			Description: "List all apps",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/deployments",
				Output:   g.Paged([]string{deployment, deployment2}, []string{deployment3, deployment4}),
				Status:   http.StatusOK},
			Expected: g.Array(deployment, deployment2, deployment3, deployment4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Deployments.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "Update deployment",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/deployments/2b56dc7b-2a14-49ea-be29-ca182b14a998",
				Output:   g.Single(deployment),
				PostForm: `{ "metadata": { "labels": { "key": "value" }, "annotations": {"note": "detailed information"}}}`,
				Status:   http.StatusOK,
			},
			Expected: deployment,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.DeploymentUpdate{
					Metadata: resource.NewMetadata().
						WithLabel("", "key", "value").
						WithAnnotation("", "note", "detailed information"),
				}
				return c.Deployments.Update(context.Background(), "2b56dc7b-2a14-49ea-be29-ca182b14a998", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
