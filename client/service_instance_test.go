package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"net/http"
	"testing"
)

func TestServiceInstances(t *testing.T) {
	g := test.NewObjectJSONGenerator(156)
	si := g.ServiceInstance()
	si2 := g.ServiceInstance()

	tests := []RouteTest{
		{
			Description: "Create service instance",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_instances",
				Output:   []string{si},
				Status:   http.StatusCreated,
				PostForm: `{
					"type": "managed",
					"name": "my_service_instance",
					"tags": ["foo", "bar", "baz"],
					"relationships": {
						"space": {
							"data": {
								"guid": "7304bc3c-7010-11ea-8840-48bf6bec2d78"
							}
						},
						"service_plan": {
							"data": {
								"guid": "e0e4417c-74ee-11ea-a604-48bf6bec2d78"
							}
						}
					}
				}`,
			},
			Expected: si,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServiceInstanceCreateManaged("my_service_instance",
					"7304bc3c-7010-11ea-8840-48bf6bec2d78", "e0e4417c-74ee-11ea-a604-48bf6bec2d78")
				r.Tags = []string{"foo", "bar", "baz"}
				return c.ServiceInstances.Create(r)
			},
		},
		{
			Description: "Delete service instance",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.ServiceInstances.Delete("62a3c0fe-5751-4f8f-97c4-28de85962ef8")
			},
		},
		{
			Description: "Get service instance",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8",
				Output:   []string{si},
				Status:   http.StatusOK},
			Expected: si,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceInstances.Get("62a3c0fe-5751-4f8f-97c4-28de85962ef8")
			},
		},
		{
			Description: "List all service instances",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_instances",
				Output:   g.Paged([]string{si}, []string{si2}),
				Status:   http.StatusOK},
			Expected: g.Array(si, si2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceInstances.ListAll(nil)
			},
		},
	}
	executeTests(tests, t)
}
