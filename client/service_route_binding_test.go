package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestServiceRouteBindings(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(324)
	svcRouteBinding := g.ServiceRouteBinding().JSON
	svcRouteBinding2 := g.ServiceRouteBinding().JSON
	svcRouteBinding3 := g.ServiceRouteBinding().JSON
	svcRouteBinding4 := g.ServiceRouteBinding().JSON
	route := g.Route().JSON
	si := g.ServiceInstance().JSON

	tests := []RouteTest{
		{
			Description: "Create route binding to managed service instance",
			Route: testutil.MockRoute{
				Method:           "POST",
				Endpoint:         "/v3/service_route_bindings",
				Status:           http.StatusCreated,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
				PostForm: `{
					"relationships": {
					  "route": {
						"data": {
						  "guid": "7304bc3c-7010-11ea-8840-48bf6bec2d78"
						}
					  },
					  "service_instance": {
						"data": {
						  "guid": "e0e4417c-74ee-11ea-a604-48bf6bec2d78"
						}
					  }
					}
				  }`,
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				r := resource.NewServiceRouteBindingCreate("7304bc3c-7010-11ea-8840-48bf6bec2d78",
					"e0e4417c-74ee-11ea-a604-48bf6bec2d78")
				return c.ServiceRouteBindings.Create(context.Background(), r)
			},
		},
		{
			Description: "Create route binding to user provided service instance",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_route_bindings",
				Output:   g.Single(svcRouteBinding),
				Status:   http.StatusCreated,
				PostForm: `{
					"relationships": {
					  "route": {
						"data": {
						  "guid": "7304bc3c-7010-11ea-8840-48bf6bec2d78"
						}
					  },
					  "service_instance": {
						"data": {
						  "guid": "e0e4417c-74ee-11ea-a604-48bf6bec2d78"
						}
					  }
					}
				  }`,
			},
			Expected:  "",
			Expected2: svcRouteBinding,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				r := resource.NewServiceRouteBindingCreate("7304bc3c-7010-11ea-8840-48bf6bec2d78",
					"e0e4417c-74ee-11ea-a604-48bf6bec2d78")
				return c.ServiceRouteBindings.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete user provided service instance route binding",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Status:   http.StatusNoContent,
			},
			Expected: "",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceRouteBindings.Delete(context.Background(), "3458647f-8358-4427-9a64-9f90392b02f7")
			},
		},
		{
			Description: "Delete managed service instance route binding",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceRouteBindings.Delete(context.Background(), "3458647f-8358-4427-9a64-9f90392b02f7")
			},
		},
		{
			Description: "Get service route binding",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Output:   g.Single(svcRouteBinding),
				Status:   http.StatusOK},
			Expected: svcRouteBinding,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceRouteBindings.Get(context.Background(), "3458647f-8358-4427-9a64-9f90392b02f7")
			},
		},
		{
			Description: "Get service route binding include route",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource: svcRouteBinding,
					Routes:   []string{route},
				}),
				Status: http.StatusOK,
			},
			Expected:  svcRouteBinding,
			Expected2: route,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceRouteBindings.GetIncludeRoute(context.Background(), "3458647f-8358-4427-9a64-9f90392b02f7")
			},
		},
		{
			Description: "Get service route binding include service instance",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource:         svcRouteBinding,
					ServiceInstances: []string{si},
				}),
				Status: http.StatusOK,
			},
			Expected:  svcRouteBinding,
			Expected2: si,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceRouteBindings.GetIncludeServiceInstance(context.Background(), "3458647f-8358-4427-9a64-9f90392b02f7")
			},
		},
		{
			Description: "Get service route binding parameters",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7/parameters",
				Output: []string{`
					{
					  "foo": "bar",
					  "foz": "baz"
					}`},
				Status: http.StatusOK,
			},
			Expected: `{ "foo": "bar", "foz": "baz" }`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceRouteBindings.GetParameters(context.Background(), "3458647f-8358-4427-9a64-9f90392b02f7")
			},
		},
		{
			Description: "List all service route bindings",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings",
				Output:   g.Paged([]string{svcRouteBinding}, []string{svcRouteBinding2}),
				Status:   http.StatusOK},
			Expected: g.Array(svcRouteBinding, svcRouteBinding2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceRouteBindings.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List all service route bindings include routes",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{svcRouteBinding, svcRouteBinding2},
						Routes:    []string{route},
					},
					testutil.PagedResult{
						Resources: []string{svcRouteBinding3, svcRouteBinding4},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(svcRouteBinding, svcRouteBinding2, svcRouteBinding3, svcRouteBinding4),
			Expected2: g.Array(route),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceRouteBindings.ListIncludeRoutesAll(context.Background(), nil)
			},
		},
		{
			Description: "List all service route bindings include service instances",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources:        []string{svcRouteBinding, svcRouteBinding2},
						ServiceInstances: []string{si},
					},
					testutil.PagedResult{
						Resources: []string{svcRouteBinding3, svcRouteBinding4},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(svcRouteBinding, svcRouteBinding2, svcRouteBinding3, svcRouteBinding4),
			Expected2: g.Array(si),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceRouteBindings.ListIncludeServiceInstancesAll(context.Background(), nil)
			},
		},
		{
			Description: "Update service route binding",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Output:   g.Single(svcRouteBinding),
				Status:   http.StatusOK,
				PostForm: `{
					"metadata": {
					  "labels": {"key": "value"},
					  "annotations": {"note": "detailed information"}
					}
				  }`,
			},
			Expected: svcRouteBinding,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.ServiceRouteBindingUpdate{
					Metadata: resource.NewMetadata().
						WithLabel("", "key", "value").
						WithAnnotation("", "note", "detailed information"),
				}
				return c.ServiceRouteBindings.Update(context.Background(), "3458647f-8358-4427-9a64-9f90392b02f7", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
