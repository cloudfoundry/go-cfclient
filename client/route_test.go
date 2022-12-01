package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(123)
	route := g.Route().JSON
	route2 := g.Route().JSON
	domain := g.Domain().JSON
	space := g.Space().JSON
	space2 := g.Space().JSON
	org := g.Organization().JSON
	routeSpaceRelationships := g.RouteSpaceRelationships().JSON
	routeDestinations := g.RouteDestinations().JSON
	routeDestinationWithLinks := g.RouteDestinationWithLinks().JSON

	tests := []RouteTest{
		{
			Description: "Create route",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/routes",
				Output:   g.Single(route),
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
				return c.Routes.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete route",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.Delete(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Delete unmapped routes for space",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/spaces/cad48e84-5a48-421e-be07-f7b4f016f581/routes",
				QueryString:      "unmapped=true",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.DeleteUnmappedRoutesForSpace(context.Background(), "cad48e84-5a48-421e-be07-f7b4f016f581")
			},
		},
		{
			Description: "Get route",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output:   g.Single(route),
				Status:   http.StatusOK,
			},
			Expected: route,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.Get(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Get route destinations",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82/destinations",
				Output:   g.Single(routeDestinations),
				Status:   http.StatusOK,
			},
			Expected: routeDestinations,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.GetDestinations(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Get shared spaces relationships",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82/relationships/shared_spaces",
				Output:   g.Single(routeSpaceRelationships),
				Status:   http.StatusOK,
			},
			Expected: routeSpaceRelationships,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.GetSharedSpacesRelationships(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Get route include domain",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource: route,
					Domains:  []string{domain},
				}),
				Status: http.StatusOK,
			},
			Expected:  route,
			Expected2: domain,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Routes.GetIncludeDomain(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Get route include space",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource: route,
					Spaces:   []string{space},
				}),
				Status: http.StatusOK,
			},
			Expected:  route,
			Expected2: space,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Routes.GetIncludeSpace(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Get route include space and organization",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource:      route,
					Spaces:        []string{space},
					Organizations: []string{org},
				}),
				Status: http.StatusOK,
			},
			Expected:  route,
			Expected2: space,
			Expected3: org,
			Action3: func(c *Client, t *testing.T) (any, any, any, error) {
				return c.Routes.GetIncludeSpaceAndOrganization(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82")
			},
		},
		{
			Description: "Check if the route is reserved",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/domains/f666ffc5-106e-4fda-b56f-568b5cf3ae9f/route_reservations",
				Output:   g.Single(`{ "matching_route": true }`),
				Status:   http.StatusOK,
			},
			Expected: "true",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.IsRouteReserved(context.Background(), "f666ffc5-106e-4fda-b56f-568b5cf3ae9f", nil)
			},
		},
		{
			Description: "Insert route destinations",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82/destinations",
				Output:   g.Single(routeDestinations),
				Status:   http.StatusCreated,
				PostForm: `{
					"destinations": [
					  {
						"app": {
						  "guid": "1cb006ee-fb05-47e1-b541-c34179ddc446"
						}
					  },
					  {
						"app": {
						  "guid": "01856e12-8ee8-11e9-98a5-bb397dbc818f",
						  "process": {
							"type": "api"
						  }
						},
						"port": 9000,
						"protocol": "http1"
					  }
					]
				  }`,
			},
			Expected: routeDestinations,
			Action: func(c *Client, t *testing.T) (any, error) {
				d := []*resource.RouteDestinationInsertOrReplace{
					resource.NewRouteDestinationInsertOrReplace("1cb006ee-fb05-47e1-b541-c34179ddc446"),
					resource.NewRouteDestinationInsertOrReplace("01856e12-8ee8-11e9-98a5-bb397dbc818f").
						WithPort(9000).WithProtocol("http1").WithProcessType("api"),
				}
				return c.Routes.InsertDestinations(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82", d)
			},
		},
		{
			Description: "List all routes",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes",
				Output:   g.Paged([]string{route}, []string{route2}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(route, route2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List all routes for an app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/758c78dc-60bc-4f84-999b-247bdc2c37fe/routes",
				Output:   g.Paged([]string{route}, []string{route2}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(route, route2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.ListForAppAll(context.Background(), "758c78dc-60bc-4f84-999b-247bdc2c37fe", nil)
			},
		},
		{
			Description: "List all routes and include domains",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{route},
						Domains:   []string{domain},
					},
					testutil.PagedResult{
						Resources: []string{route2},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(route, route2),
			Expected2: g.Array(domain),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Routes.ListIncludeDomainsAll(context.Background(), nil)
			},
		},
		{
			Description: "List all routes and include spaces",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{route},
						Spaces:    []string{space},
					},
					testutil.PagedResult{
						Resources: []string{route2},
						Spaces:    []string{space2},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(route, route2),
			Expected2: g.Array(space, space2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Routes.ListIncludeSpacesAll(context.Background(), nil)
			},
		},
		{
			Description: "List all routes and include spaces and organizations",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/routes",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources:     []string{route},
						Spaces:        []string{space},
						Organizations: []string{org},
					},
					testutil.PagedResult{
						Resources: []string{route2},
						Spaces:    []string{space2},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(route, route2),
			Expected2: g.Array(space, space2),
			Expected3: g.Array(org),
			Action3: func(c *Client, t *testing.T) (any, any, any, error) {
				return c.Routes.ListIncludeSpacesAndOrganizationsAll(context.Background(), nil)
			},
		},
		{
			Description: "Remove destination for route",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82/destinations/1cb006ee-fb05-47e1-b541-c34179ddc446",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Routes.RemoveDestination(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82",
					"1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "Replace route destinations",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82/destinations",
				Output:   g.Single(routeDestinations),
				Status:   http.StatusOK,
				PostForm: `{
					"destinations": [
					  {
						"app": {
						  "guid": "1cb006ee-fb05-47e1-b541-c34179ddc446"
						},
 						"weight": 61
					  },
					  {
						"app": {
						  "guid": "01856e12-8ee8-11e9-98a5-bb397dbc818f",
						  "process": {
							"type": "api"
						  }
						},
						"weight": 39,
						"port": 9000,
						"protocol": "http1"
					  }
					]
				  }`,
			},
			Expected: routeDestinations,
			Action: func(c *Client, t *testing.T) (any, error) {
				d := []*resource.RouteDestinationInsertOrReplace{
					resource.NewRouteDestinationInsertOrReplace("1cb006ee-fb05-47e1-b541-c34179ddc446").
						WithWeight(61),
					resource.NewRouteDestinationInsertOrReplace("01856e12-8ee8-11e9-98a5-bb397dbc818f").
						WithPort(9000).
						WithProtocol("http1").
						WithProcessType("api").
						WithWeight(39),
				}
				return c.Routes.ReplaceDestinations(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82", d)
			},
		},
		{
			Description: "Share route with space",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82/relationships/shared_spaces",
				Output:   g.Single(routeSpaceRelationships),
				Status:   http.StatusOK,
				PostForm: `{ "data": [{ "guid":"68d54d31-9b3a-463b-ba94-e8e4c32edbac" }]}`,
			},
			Expected: routeSpaceRelationships,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.ShareWithSpace(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82", "68d54d31-9b3a-463b-ba94-e8e4c32edbac")
			},
		},
		{
			Description: "Un-Share route with spaces",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82/relationships/shared_spaces/68d54d31-9b3a-463b-ba94-e8e4c32edbac",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Routes.UnShareWithSpaces(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82", []string{"68d54d31-9b3a-463b-ba94-e8e4c32edbac"})
			},
		},
		{
			Description: "Transfer route ownership",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82/relationships/space",
				Status:   http.StatusNoContent,
				PostForm: `{ "data": { "guid": "68d54d31-9b3a-463b-ba94-e8e4c32edbac"} }`,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Routes.TransferOwnership(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82", "68d54d31-9b3a-463b-ba94-e8e4c32edbac")
			},
		},
		{
			Description: "Update destination protocol",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82/destinations/6e2123df-db4f-4d89-941b-5c79a0b0aa4a",
				Output:   g.Single(routeDestinationWithLinks),
				Status:   http.StatusOK,
				PostForm: `{"protocol": "http2"}`,
			},
			Expected: routeDestinationWithLinks,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Routes.UpdateDestinationProtocol(context.Background(),
					"5a85c020-3e3d-42a5-a475-5084c5357e82",
					"6e2123df-db4f-4d89-941b-5c79a0b0aa4a",
					"http2")
			},
		},
		{
			Description: "Update route",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/routes/5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output:   g.Single(route),
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": {"key": "value"}, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: route,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.RouteUpdate{
					Metadata: resource.NewMetadata().
						WithLabel("", "key", "value").
						WithAnnotation("", "note", "detailed information"),
				}
				return c.Routes.Update(context.Background(), "5a85c020-3e3d-42a5-a475-5084c5357e82", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
