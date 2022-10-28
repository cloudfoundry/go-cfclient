package client

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDroplets(t *testing.T) {
	g := test.NewObjectJSONGenerator(2)
	droplet := g.Droplet()
	droplet2 := g.Droplet()
	droplet3 := g.Droplet()
	droplet4 := g.Droplet()
	dropletAssociation := g.DropletAssociation()

	tests := []RouteTest{
		{
			Description: "Set current droplet association for app",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/apps/bf75e72f-f1ed-4815-9e28-048595a35b6c/relationships/current_droplet",
				Output:   []string{dropletAssociation},
				Status:   http.StatusOK,
				PostForm: `{"data":{"guid":"3fc0916f-2cea-4f3a-ae53-048388baa6bd"}}`,
			},
			Expected: dropletAssociation,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.SetCurrentAssociationForApp("bf75e72f-f1ed-4815-9e28-048595a35b6c", "3fc0916f-2cea-4f3a-ae53-048388baa6bd")
			},
		},
		{
			Description: "Get current droplet association for app",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/bf75e72f-f1ed-4815-9e28-048595a35b6c/relationships/current_droplet",
				Output:   []string{dropletAssociation},
				Status:   http.StatusOK,
			},
			Expected: dropletAssociation,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.GetCurrentAssociationForApp("bf75e72f-f1ed-4815-9e28-048595a35b6c")
			},
		},
		{
			Description: "Get current droplet for app",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/bf75e72f-f1ed-4815-9e28-048595a35b6c/droplets/current",
				Output:   []string{droplet},
				Status:   http.StatusOK,
			},
			Expected: droplet,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.GetCurrentForApp("bf75e72f-f1ed-4815-9e28-048595a35b6c")
			},
		},
		{
			Description: "Get droplet",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4",
				Output:   []string{droplet},
				Status:   http.StatusOK,
			},
			Expected: droplet,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.Get("59c3d133-2b83-46f3-960e-7765a129aea4")
			},
		},
		{
			Description: "List first page of droplets",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/droplets",
				Output:   g.Paged("droplets", []string{droplet}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(droplet),
			Action: func(c *Client, t *testing.T) (any, error) {
				ds, _, err := c.Droplets.List(NewDropletListOptions())
				return ds, err
			},
		},
		{
			Description: "List all droplets",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/droplets",
				Output:   g.Paged("droplets", []string{droplet, droplet2}, []string{droplet3, droplet4}),
				Status:   http.StatusOK},
			Expected: g.Array(droplet, droplet2, droplet3, droplet4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.ListAll(nil)
			},
		},
		{
			Description: "List first page of droplets for a package",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/packages/8222f76a-9e09-4360-b3aa-1ed329945e92/droplets",
				Output:   g.Paged("droplets", []string{droplet}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(droplet),
			Action: func(c *Client, t *testing.T) (any, error) {
				ds, _, err := c.Droplets.ListForPackage("8222f76a-9e09-4360-b3aa-1ed329945e92", NewDropletPackageListOptions())
				return ds, err
			},
		},
		{
			Description: "List first page of droplets for an app",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/bf75e72f-f1ed-4815-9e28-048595a35b6c/droplets",
				Output:   g.Paged("droplets", []string{droplet}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(droplet),
			Action: func(c *Client, t *testing.T) (any, error) {
				ds, _, err := c.Droplets.ListForApp("bf75e72f-f1ed-4815-9e28-048595a35b6c", NewDropletAppListOptions())
				return ds, err
			},
		},
		{
			Description: "Create droplet",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/droplets",
				Output:   []string{droplet},
				Status:   http.StatusCreated,
			},
			Expected: droplet,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewDropletCreate("bf75e72f-f1ed-4815-9e28-048595a35b6c")
				return c.Droplets.Create(r)
			},
		},
		{
			Description: "Delete droplet",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Droplets.Delete("59c3d133-2b83-46f3-960e-7765a129aea4")
			},
		},
		{
			Description: "Update droplet",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4",
				Output:   []string{droplet},
				Status:   http.StatusOK,
			},
			Expected: droplet,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.Update("59c3d133-2b83-46f3-960e-7765a129aea4", &resource.DropletUpdate{})
			},
		},
		{
			Description: "Copy droplet",
			Route: MockRoute{
				Method:      "POST",
				Endpoint:    "/v3/droplets",
				QueryString: "source_guid=59c3d133-2b83-46f3-960e-7765a129aea4",
				Output:      []string{droplet},
				Status:      http.StatusCreated,
				PostForm:    `{ "relationships": { "app": { "data": { "guid": "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b" }}}}`,
			},
			Expected: droplet,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Droplets.Copy("59c3d133-2b83-46f3-960e-7765a129aea4", "8d1f1d2e-08b1-4a10-a8df-471a1418cb8b")
			},
		},
		{
			Description: "Download droplet",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/droplets/59c3d133-2b83-46f3-960e-7765a129aea4/download",
				Output:   []string{"droplet bits..."},
				Status:   http.StatusOK,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				reader, err := c.Droplets.Download("59c3d133-2b83-46f3-960e-7765a129aea4")
				require.NoError(t, err)
				buf := new(strings.Builder)
				_, err = io.Copy(buf, reader)
				require.NoError(t, err)
				require.Equal(t, "droplet bits...", buf.String())
				return nil, nil
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
