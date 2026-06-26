package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/cloudfoundry/go-cfclient/v3/testutil"
)

func TestRoutePolicies(t *testing.T) {
	g := testutil.NewObjectJSONGenerator()
	routePolicy := g.RoutePolicy().JSON
	routePolicy2 := g.RoutePolicy().JSON
	route1 := g.Route().JSON
	route2 := g.Route().JSON
	app1 := g.Application().JSON
	app2 := g.Application().JSON

	tests := []RouteTest{
		{
			Description: "Create route policy",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/route_policies",
				Output:   g.Single(routePolicy),
				Status:   http.StatusCreated,
				PostForm: `{
					"source": "1cb006ee-fb05-47e1-b541-c34179ddc446",
					"relationships": {
						"route": {
							"data": { "guid": "5a85c020-3e3d-42a5-a475-5084c5357e82" }
						},
						"app": {
							"data": { "guid": "1cb006ee-fb05-47e1-b541-c34179ddc446" }
						}
					}
				}`,
			},
			Expected: routePolicy,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.RoutePolicyCreate{
					Source: "1cb006ee-fb05-47e1-b541-c34179ddc446",
					Relationships: resource.RoutePolicyRelationships{
						Route: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "5a85c020-3e3d-42a5-a475-5084c5357e82",
							},
						},
						App: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "1cb006ee-fb05-47e1-b541-c34179ddc446",
							},
						},
					},
				}
				return c.RoutePolicies.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete route policy",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/route_policies/c8dcf27f-39a3-466a-9cbf-1d3c31a43b93",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.RoutePolicies.Delete(context.Background(), "c8dcf27f-39a3-466a-9cbf-1d3c31a43b93")
			},
		},
		{
			Description: "Get route policy",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies/c8dcf27f-39a3-466a-9cbf-1d3c31a43b93",
				Output:   g.Single(routePolicy),
				Status:   http.StatusOK,
			},
			Expected: routePolicy,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.RoutePolicies.Get(context.Background(), "c8dcf27f-39a3-466a-9cbf-1d3c31a43b93")
			},
		},
		{
			Description: "Get route policy with include route",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies/c8dcf27f-39a3-466a-9cbf-1d3c31a43b93",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{routePolicy},
						Routes:    []string{route1},
					}),
				Status: http.StatusOK,
			},
			Expected:  routePolicy,
			Expected2: route1,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.RoutePolicies.GetIncludeRoute(context.Background(), "c8dcf27f-39a3-466a-9cbf-1d3c31a43b93")
			},
		},
		{
			Description: "Get route policy with include source",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies/c8dcf27f-39a3-466a-9cbf-1d3c31a43b93",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{routePolicy},
						Apps:      []string{app1},
					}),
				Status: http.StatusOK,
			},
			Expected:  routePolicy,
			Expected2: app1,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.RoutePolicies.GetIncludeSource(context.Background(), "c8dcf27f-39a3-466a-9cbf-1d3c31a43b93")
			},
		},
		{
			Description: "List first page of route policies",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies",
				Output:   g.SinglePaged(routePolicy),
				Status:   http.StatusOK,
			},
			Expected: g.Array(routePolicy),
			Action: func(c *Client, t *testing.T) (any, error) {
				policies, _, err := c.RoutePolicies.List(context.Background(), NewRoutePolicyListOptions())
				return policies, err
			},
		},
		{
			Description: "List all route policies",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies",
				Output:   g.Paged([]string{routePolicy}, []string{routePolicy2}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(routePolicy, routePolicy2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.RoutePolicies.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "First route policy with no matches",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies",
				Output:   g.Paged([]string{}),
				Status:   http.StatusOK,
			},
			Expected: "error",
			Action: func(c *Client, t *testing.T) (any, error) {
				_, err := c.RoutePolicies.First(context.Background(), NewRoutePolicyListOptions())
				if err != nil {
					return "error", nil
				}
				return "no error", nil
			},
		},
		{
			Description: "First route policy",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies",
				Output:   g.SinglePaged(routePolicy),
				Status:   http.StatusOK,
			},
			Expected: routePolicy,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.RoutePolicies.First(context.Background(), NewRoutePolicyListOptions())
			},
		},
		{
			Description: "Single route policy with no matches",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies",
				Output:   g.Paged([]string{}),
				Status:   http.StatusOK,
			},
			Expected: "error",
			Action: func(c *Client, t *testing.T) (any, error) {
				_, err := c.RoutePolicies.Single(context.Background(), NewRoutePolicyListOptions())
				if err != nil {
					return "error", nil
				}
				return "no error", nil
			},
		},
		{
			Description: "Single route policy",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies",
				Output:   g.SinglePaged(routePolicy),
				Status:   http.StatusOK,
			},
			Expected: routePolicy,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.RoutePolicies.Single(context.Background(), NewRoutePolicyListOptions())
			},
		},
		{
			Description: "List route policies with filter by route GUID",
			Route: testutil.MockRoute{
				Method:      "GET",
				Endpoint:    "/v3/route_policies",
				QueryString: "page=1&per_page=50&route_guids=5a85c020-3e3d-42a5-a475-5084c5357e82",
				Output:      g.SinglePaged(routePolicy),
				Status:      http.StatusOK,
			},
			Expected: g.Array(routePolicy),
			Action: func(c *Client, t *testing.T) (any, error) {
				opts := NewRoutePolicyListOptions()
				opts.RouteGUIDs = Filter{
					Values: []string{"5a85c020-3e3d-42a5-a475-5084c5357e82"},
				}
				policies, _, err := c.RoutePolicies.List(context.Background(), opts)
				return policies, err
			},
		},
		{
			Description: "List route policies with filter by space GUID",
			Route: testutil.MockRoute{
				Method:      "GET",
				Endpoint:    "/v3/route_policies",
				QueryString: "page=1&per_page=50&space_guids=33d27af8-788d-4de5-8f37-fb80d517f2ed",
				Output:      g.SinglePaged(routePolicy),
				Status:      http.StatusOK,
			},
			Expected: g.Array(routePolicy),
			Action: func(c *Client, t *testing.T) (any, error) {
				opts := NewRoutePolicyListOptions()
				opts.SpaceGUIDs = Filter{
					Values: []string{"33d27af8-788d-4de5-8f37-fb80d517f2ed"},
				}
				policies, _, err := c.RoutePolicies.List(context.Background(), opts)
				return policies, err
			},
		},
		{
			Description: "List route policies with filter by source GUID",
			Route: testutil.MockRoute{
				Method:      "GET",
				Endpoint:    "/v3/route_policies",
				QueryString: "page=1&per_page=50&source_guids=1cb006ee-fb05-47e1-b541-c34179ddc446",
				Output:      g.SinglePaged(routePolicy),
				Status:      http.StatusOK,
			},
			Expected: g.Array(routePolicy),
			Action: func(c *Client, t *testing.T) (any, error) {
				opts := NewRoutePolicyListOptions()
				opts.SourceGUIDs = Filter{
					Values: []string{"1cb006ee-fb05-47e1-b541-c34179ddc446"},
				}
				policies, _, err := c.RoutePolicies.List(context.Background(), opts)
				return policies, err
			},
		},
		{
			Description: "List route policies with filter by source",
			Route: testutil.MockRoute{
				Method:      "GET",
				Endpoint:    "/v3/route_policies",
				QueryString: "page=1&per_page=50&sources=network_policy",
				Output:      g.SinglePaged(routePolicy),
				Status:      http.StatusOK,
			},
			Expected: g.Array(routePolicy),
			Action: func(c *Client, t *testing.T) (any, error) {
				opts := NewRoutePolicyListOptions()
				opts.Sources = Filter{
					Values: []string{"network_policy"},
				}
				policies, _, err := c.RoutePolicies.List(context.Background(), opts)
				return policies, err
			},
		},
		{
			Description: "List route policies with filter by organization GUID",
			Route: testutil.MockRoute{
				Method:      "GET",
				Endpoint:    "/v3/route_policies",
				QueryString: "organization_guids=3a5f687b-2ce8-4ade-be75-8eca99b0db8b&page=1&per_page=50",
				Output:      g.SinglePaged(routePolicy),
				Status:      http.StatusOK,
			},
			Expected: g.Array(routePolicy),
			Action: func(c *Client, t *testing.T) (any, error) {
				opts := NewRoutePolicyListOptions()
				opts.OrganizationGUIDs = Filter{
					Values: []string{"3a5f687b-2ce8-4ade-be75-8eca99b0db8b"},
				}
				policies, _, err := c.RoutePolicies.List(context.Background(), opts)
				return policies, err
			},
		},
		{
			Description: "List route policies with include route",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies",
				Output:   g.SinglePaged(routePolicy),
				Status:   http.StatusOK,
			},
			Expected: g.Array(routePolicy),
			Action: func(c *Client, t *testing.T) (any, error) {
				policies, _, _, err := c.RoutePolicies.ListIncludeRoute(context.Background(), NewRoutePolicyListOptions())
				return policies, err
			},
		},
		{
			Description: "List all route policies with include route",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{routePolicy, routePolicy2},
						Routes:    []string{route1, route2},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(routePolicy, routePolicy2),
			Expected2: g.Array(route1, route2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.RoutePolicies.ListIncludeRouteAll(context.Background(), nil)
			},
		},
		{
			Description: "List route policies with include source",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies",
				Output:   g.SinglePaged(routePolicy),
				Status:   http.StatusOK,
			},
			Expected: g.Array(routePolicy),
			Action: func(c *Client, t *testing.T) (any, error) {
				policies, _, _, _, _, err := c.RoutePolicies.ListIncludeSource(context.Background(), NewRoutePolicyListOptions())
				return policies, err
			},
		},
		{
			Description: "List all route policies with include source",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/route_policies",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{routePolicy, routePolicy2},
						Apps:      []string{app1, app2},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(routePolicy, routePolicy2),
			Expected2: g.Array(app1, app2),
			Action4: func(c *Client, t *testing.T) (any, any, any, any, error) {
				return c.RoutePolicies.ListIncludeSourceAll(context.Background(), nil)
			},
		},
		{
			Description: "Update route policy",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/route_policies/c8dcf27f-39a3-466a-9cbf-1d3c31a43b93",
				Output:   g.Single(routePolicy),
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": {"key": "value"}, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: routePolicy,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.RoutePolicyUpdate{
					Metadata: resource.NewMetadata().
						WithLabel("", "key", "value").
						WithAnnotation("", "note", "detailed information"),
				}
				return c.RoutePolicies.Update(context.Background(), "c8dcf27f-39a3-466a-9cbf-1d3c31a43b93", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
