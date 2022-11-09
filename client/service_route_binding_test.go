package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"net/http"
	"testing"
)

func TestServiceRouteBindings(t *testing.T) {
	g := test.NewObjectJSONGenerator(324)
	svcRouteBinding := g.ServiceRouteBinding()
	svcRouteBinding2 := g.ServiceRouteBinding()
	svcRouteBinding3 := g.ServiceRouteBinding()
	svcRouteBinding4 := g.ServiceRouteBinding()
	route := g.Route()
	si := g.ServiceInstance()

	tests := []RouteTest{
		{
			Description: "Create route binding",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_route_bindings",
				Output:   []string{svcRouteBinding},
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
			Expected: svcRouteBinding,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServiceRouteBindingCreate("7304bc3c-7010-11ea-8840-48bf6bec2d78",
					"e0e4417c-74ee-11ea-a604-48bf6bec2d78")
				return c.ServiceRouteBindings.Create(r)
			},
		},
		{
			Description: "Delete service plan",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.ServiceRouteBindings.Delete("3458647f-8358-4427-9a64-9f90392b02f7")
				return nil, err
			},
		},
		{
			Description: "Get service route binding",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Output:   []string{svcRouteBinding},
				Status:   http.StatusOK},
			Expected: svcRouteBinding,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceRouteBindings.Get("3458647f-8358-4427-9a64-9f90392b02f7")
			},
		},
		{
			Description: "Get service route binding include route",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Output: g.ResourceWithInclude(test.ResourceResult{
					Resource: svcRouteBinding,
					Routes:   []string{route},
				}),
				Status: http.StatusOK,
			},
			Expected:  svcRouteBinding,
			Expected2: route,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceRouteBindings.GetIncludeRoute("3458647f-8358-4427-9a64-9f90392b02f7")
			},
		},
		{
			Description: "Get service route binding include service instance",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Output: g.ResourceWithInclude(test.ResourceResult{
					Resource:         svcRouteBinding,
					ServiceInstances: []string{si},
				}),
				Status: http.StatusOK,
			},
			Expected:  svcRouteBinding,
			Expected2: si,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceRouteBindings.GetIncludeServiceInstance("3458647f-8358-4427-9a64-9f90392b02f7")
			},
		},
		{
			Description: "Get service route binding parameters",
			Route: MockRoute{
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
				return c.ServiceRouteBindings.GetParameters("3458647f-8358-4427-9a64-9f90392b02f7")
			},
		},
		{
			Description: "List all service route bindings",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings",
				Output:   g.Paged([]string{svcRouteBinding}, []string{svcRouteBinding2}),
				Status:   http.StatusOK},
			Expected: g.Array(svcRouteBinding, svcRouteBinding2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceRouteBindings.ListAll(nil)
			},
		},
		{
			Description: "List all service route bindings include routes",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources: []string{svcRouteBinding, svcRouteBinding2},
						Routes:    []string{route},
					},
					test.PagedResult{
						Resources: []string{svcRouteBinding3, svcRouteBinding4},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(svcRouteBinding, svcRouteBinding2, svcRouteBinding3, svcRouteBinding4),
			Expected2: g.Array(route),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceRouteBindings.ListIncludeRoutesAll(nil)
			},
		},
		{
			Description: "List all service route bindings include service instances",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_route_bindings",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources:        []string{svcRouteBinding, svcRouteBinding2},
						ServiceInstances: []string{si},
					},
					test.PagedResult{
						Resources: []string{svcRouteBinding3, svcRouteBinding4},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(svcRouteBinding, svcRouteBinding2, svcRouteBinding3, svcRouteBinding4),
			Expected2: g.Array(si),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceRouteBindings.ListIncludeServiceInstancesAll(nil)
			},
		},
		{
			Description: "Update service route binding",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_route_bindings/3458647f-8358-4427-9a64-9f90392b02f7",
				Output:   []string{svcRouteBinding},
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
					Metadata: resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.ServiceRouteBindings.Update("3458647f-8358-4427-9a64-9f90392b02f7", r)
			},
		},
	}
	executeTests(tests, t)
}
