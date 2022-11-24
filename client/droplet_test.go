package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDroplets(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(2)
	droplet := g.Droplet().JSON
	droplet2 := g.Droplet().JSON
	droplet3 := g.Droplet().JSON
	droplet4 := g.Droplet().JSON
	dropletAssociation := g.DropletAssociation().JSON

	blobstoreServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("droplet bits..."))
	}))
	defer blobstoreServer.Close()

	tests := []RouteTest{
		{
			Description: "Set current droplet association for app",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/apps/bf75e72f-f1ed-4815-9e28-048595a35b6c/relationships/current_droplet",
				Output:   g.Single(dropletAssociation),
				Status:   http.StatusOK,
				PostForm: `{"data":{"guid":"3fc0916f-2cea-4f3a-ae53-048388baa6bd"}}`,
			},
			Expected: dropletAssociation,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.SetCurrentAssociationForApp(context.Background(), "bf75e72f-f1ed-4815-9e28-048595a35b6c", "3fc0916f-2cea-4f3a-ae53-048388baa6bd")
			},
		},
		{
			Description: "Get current droplet association for app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/bf75e72f-f1ed-4815-9e28-048595a35b6c/relationships/current_droplet",
				Output:   g.Single(dropletAssociation),
				Status:   http.StatusOK,
			},
			Expected: dropletAssociation,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.GetCurrentAssociationForApp(context.Background(), "bf75e72f-f1ed-4815-9e28-048595a35b6c")
			},
		},
		{
			Description: "Get current droplet for app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/bf75e72f-f1ed-4815-9e28-048595a35b6c/droplets/current",
				Output:   g.Single(droplet),
				Status:   http.StatusOK,
			},
			Expected: droplet,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.GetCurrentForApp(context.Background(), "bf75e72f-f1ed-4815-9e28-048595a35b6c")
			},
		},
		{
			Description: "Get droplet",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4",
				Output:   g.Single(droplet),
				Status:   http.StatusOK,
			},
			Expected: droplet,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.Get(context.Background(), "59c3d133-2b83-46f3-960e-7765a129aea4")
			},
		},
		{
			Description: "List first page of droplets",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/droplets",
				Output:   g.SinglePaged(droplet),
				Status:   http.StatusOK,
			},
			Expected: g.Array(droplet),
			Action: func(c *Client, t *testing.T) (any, error) {
				ds, _, err := c.Droplets.List(context.Background(), NewDropletListOptions())
				return ds, err
			},
		},
		{
			Description: "List all droplets",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/droplets",
				Output:   g.Paged([]string{droplet, droplet2}, []string{droplet3, droplet4}),
				Status:   http.StatusOK},
			Expected: g.Array(droplet, droplet2, droplet3, droplet4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List first page of droplets for a package",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/packages/8222f76a-9e09-4360-b3aa-1ed329945e92/droplets",
				Output:   g.SinglePaged(droplet),
				Status:   http.StatusOK,
			},
			Expected: g.Array(droplet),
			Action: func(c *Client, t *testing.T) (any, error) {
				ds, _, err := c.Droplets.ListForPackage(context.Background(), "8222f76a-9e09-4360-b3aa-1ed329945e92", NewDropletPackageListOptions())
				return ds, err
			},
		},
		{
			Description: "List all droplets for a package",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/packages/8222f76a-9e09-4360-b3aa-1ed329945e92/droplets",
				Output:   g.Paged([]string{droplet, droplet2}, []string{droplet3}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(droplet, droplet2, droplet3),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.ListForPackageAll(context.Background(), "8222f76a-9e09-4360-b3aa-1ed329945e92", nil)
			},
		},
		{
			Description: "List first page of droplets for an app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/bf75e72f-f1ed-4815-9e28-048595a35b6c/droplets",
				Output:   g.SinglePaged(droplet),
				Status:   http.StatusOK,
			},
			Expected: g.Array(droplet),
			Action: func(c *Client, t *testing.T) (any, error) {
				ds, _, err := c.Droplets.ListForApp(context.Background(), "bf75e72f-f1ed-4815-9e28-048595a35b6c", NewDropletAppListOptions())
				return ds, err
			},
		},
		{
			Description: "Create droplet",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/droplets",
				Output:   g.Single(droplet),
				Status:   http.StatusCreated,
			},
			Expected: droplet,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewDropletCreate("bf75e72f-f1ed-4815-9e28-048595a35b6c")
				return c.Droplets.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete droplet",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.Delete(context.Background(), "59c3d133-2b83-46f3-960e-7765a129aea4")
			},
		},
		{
			Description: "Update droplet",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4",
				Output:   g.Single(droplet),
				Status:   http.StatusOK,
			},
			Expected: droplet,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.Update(context.Background(), "59c3d133-2b83-46f3-960e-7765a129aea4", &resource.DropletUpdate{})
			},
		},
		{
			Description: "Copy droplet",
			Route: testutil.MockRoute{
				Method:      "POST",
				Endpoint:    "/v3/droplets",
				QueryString: "source_guid=59c3d133-2b83-46f3-960e-7765a129aea4",
				Output:      g.Single(droplet),
				Status:      http.StatusCreated,
				PostForm:    `{ "relationships": { "app": { "data": { "guid": "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b" }}}}`,
			},
			Expected: droplet,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.Copy(context.Background(), "59c3d133-2b83-46f3-960e-7765a129aea4", "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b")
			},
		},
		{
			Description: "Download droplet",
			Route: testutil.MockRoute{
				Method:           "GET",
				Endpoint:         "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4/download",
				Status:           http.StatusFound,
				RedirectLocation: blobstoreServer.URL,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				reader, err := c.Droplets.Download(context.Background(), "59c3d133-2b83-46f3-960e-7765a129aea4")
				require.NoError(t, err)
				buf := new(strings.Builder)
				_, err = io.Copy(buf, reader)
				require.NoError(t, err)
				require.Equal(t, "droplet bits...", buf.String())
				return nil, nil
			},
		},
		{
			Description: "Upload droplet",
			Route: testutil.MockRoute{
				Method:           "POST",
				Endpoint:         "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4/upload",
				Output:           g.Single(droplet),
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected:  "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Expected2: droplet,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				tgzFile := strings.NewReader("droplet")
				return c.Droplets.Upload(context.Background(), "59c3d133-2b83-46f3-960e-7765a129aea4", tgzFile)
			},
		},
	}
	ExecuteTests(tests, t)
}
