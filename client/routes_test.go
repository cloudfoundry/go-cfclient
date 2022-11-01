package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	g := test.NewObjectJSONGenerator(123)
	route := g.Route()
	route2 := g.Route()

	tests := []RouteTest{
		{
			Description: "Create route",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/routes",
				Output:   []string{route},
				Status:   http.StatusCreated,
				PostForm: `{
					"host": "a-hostname",
					"path": "/some_path",
					"port": 6666,
					"relationships": {
					  "domain": {
						"data": { "guid": "a99f869d-151a-4a80-95b7-653ada640824" }
					  },
					  "space": {
						"data": { "guid": "33d27af8-788d-4de5-8f37-fb80d517f2ed" }
					  }
					}
				  }`,
			},
			Expected: route,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewRouteCreateWithHost("a99f869d-151a-4a80-95b7-653ada640824",
					"33d27af8-788d-4de5-8f37-fb80d517f2ed",
					"a-hostname",
					"/some_path",
					6666)
				return c.Routes.Create(r)
			},
		},
		{
			Description: "Delete route",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Routes.Delete("5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Get route",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output:   []string{route},
				Status:   http.StatusOK},
			Expected: route,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.Get("5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "List all routes",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes",
				Output:   g.Paged([]string{route}, []string{route2}),
				Status:   http.StatusOK},
			Expected: g.Array(route, route2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.ListAll(nil)
			},
		},
		{
			Description: "List all routes for an app",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/758c78dc-60bc-4f84-999b-247bdc2c37fe/routes",
				Output:   g.Paged([]string{route}, []string{route2}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(route, route2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.ListForAppAll("758c78dc-60bc-4f84-999b-247bdc2c37fe", nil)
			},
		},
		{
			Description: "Update route",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output:   []string{route},
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": {"key": "value"}, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: route,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.RouteUpdate{
					Metadata: &resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.Routes.Update("5a85c020-3e3d-42a5-a475-5084c5357e82", r)
			},
		},
	}
	executeTests(tests, t)
}
