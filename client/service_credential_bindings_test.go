package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"net/http"
	"testing"
)

func TestServiceCredentialBindings(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	scb := g.ServiceCredentialBinding()
	scb2 := g.ServiceCredentialBinding()
	scb3 := g.ServiceCredentialBinding()
	scb4 := g.ServiceCredentialBinding()
	app := g.Application()
	app2 := g.Application()
	app3 := g.Application()
	app4 := g.Application()
	si := g.ServiceInstance()

	tests := []RouteTest{
		{
			Description: "Create service credential binding",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_credential_bindings",
				Output:   []string{scb},
				Status:   http.StatusCreated,
				PostForm: `{
					"type": "app",
					"name": "some-binding-name",
					"relationships": {
					  "service_instance": {
						"data": {
						  "guid": "7304bc3c-7010-11ea-8840-48bf6bec2d78"
						}
					  },
					  "app": {
						"data": {
						  "guid": "e0e4417c-74ee-11ea-a604-48bf6bec2d78"
						}
					  }
					},
					"parameters": {
					  "key1": "value1",
					  "key2": "value2"
					}
				  }`,
			},
			Expected: scb,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServiceCredentialBindingCreateApp(
					"7304bc3c-7010-11ea-8840-48bf6bec2d78", "e0e4417c-74ee-11ea-a604-48bf6bec2d78").
					WithName("some-binding-name").
					WithJSONParameters(`{
					  "key1": "value1",
					  "key2": "value2"
					}`)
				return c.ServiceCredentialBindings.Create(r)
			},
		},
		{
			Description: "Delete service credential binding",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.ServiceCredentialBindings.Delete("59ba6d78-6a21-4321-83a9-f7eacd88b08d")
			},
		},
		{
			Description: "Get service credential binding",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d",
				Output:   []string{scb},
				Status:   http.StatusOK},
			Expected: scb,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceCredentialBindings.Get("59ba6d78-6a21-4321-83a9-f7eacd88b08d")
			},
		},
		{
			Description: "Get service credential binding and app",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d",
				Output: g.ResourceWithInclude(test.ResourceResult{
					Resource: scb,
					Apps:     []string{app},
				}),
				Status: http.StatusOK,
			},
			Expected:  scb,
			Expected2: app,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceCredentialBindings.GetIncludeApp("59ba6d78-6a21-4321-83a9-f7eacd88b08d")
			},
		},
		{
			Description: "Get service credential binding and service instance",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d",
				Output: g.ResourceWithInclude(test.ResourceResult{
					Resource:         scb,
					ServiceInstances: []string{si},
				}),
				Status: http.StatusOK,
			},
			Expected:  scb,
			Expected2: si,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceCredentialBindings.GetIncludeServiceInstance("59ba6d78-6a21-4321-83a9-f7eacd88b08d")
			},
		},
		{
			Description: "List all service credential bindings",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings",
				Output:   g.Paged([]string{scb}, []string{scb2}),
				Status:   http.StatusOK},
			Expected: g.Array(scb, scb2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceCredentialBindings.ListAll(nil)
			},
		},
		{
			Description: "List all service credential bindings and include apps",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources: []string{scb, scb2},
						Apps:      []string{app, app2},
					},
					test.PagedResult{
						Resources: []string{scb3, scb4},
						Apps:      []string{app3, app4},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(scb, scb2, scb3, scb4),
			Expected2: g.Array(app, app2, app3, app4),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceCredentialBindings.ListIncludeAppsAll(nil)
			},
		},
		{
			Description: "List all service credential bindings and include service instances",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources:        []string{scb, scb2},
						ServiceInstances: []string{si},
					},
					test.PagedResult{
						Resources: []string{scb3, scb4},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(scb, scb2, scb3, scb4),
			Expected2: g.Array(si),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceCredentialBindings.ListIncludeServiceInstancesAll(nil)
			},
		},
		{
			Description: "Update service credential binding",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d",
				Output:   []string{scb},
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": {"foo": "bar"}, "annotations": {"baz": "qux"} }}`,
			},
			Expected: scb,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.ServiceCredentialBindingUpdate{
					Metadata: resource.Metadata{
						Labels: map[string]string{
							"foo": "bar",
						},
						Annotations: map[string]string{
							"baz": "qux",
						},
					},
				}
				return c.ServiceCredentialBindings.Update("59ba6d78-6a21-4321-83a9-f7eacd88b08d", r)
			},
		},
	}
	executeTests(tests, t)
}
