package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestServiceCredentialBindings(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	scb := g.ServiceCredentialBinding().JSON
	scb2 := g.ServiceCredentialBinding().JSON
	scb3 := g.ServiceCredentialBinding().JSON
	scb4 := g.ServiceCredentialBinding().JSON
	scbd := g.ServiceCredentialBindingDetails().JSON
	app := g.Application().JSON
	app2 := g.Application().JSON
	app3 := g.Application().JSON
	app4 := g.Application().JSON
	si := g.ServiceInstance().JSON

	tests := []RouteTest{
		{
			Description: "Create app service credential binding for user provided service",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_credential_bindings",
				Output:   g.Single(scb),
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
				_, binding, err := c.ServiceCredentialBindings.Create(context.Background(), r)
				return binding, err
			},
		},
		{
			Description: "Create app service credential binding for managed service instance",
			Route: testutil.MockRoute{
				Method:           "POST",
				Endpoint:         "/v3/service_credential_bindings",
				Output:           g.Single(scb),
				Status:           http.StatusCreated,
				RedirectLocation: "https://api.example.org/v3/jobs/af5c57f6-8769-41fa-a499-2c84ed896788",
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
			Expected: "af5c57f6-8769-41fa-a499-2c84ed896788",
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				r := resource.NewServiceCredentialBindingCreateApp(
					"7304bc3c-7010-11ea-8840-48bf6bec2d78", "e0e4417c-74ee-11ea-a604-48bf6bec2d78").
					WithName("some-binding-name").
					WithJSONParameters(`{
					  "key1": "value1",
					  "key2": "value2"
					}`)
				return c.ServiceCredentialBindings.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete service credential binding",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.ServiceCredentialBindings.Delete(context.Background(), "59ba6d78-6a21-4321-83a9-f7eacd88b08d")
			},
		},
		{
			Description: "Get service credential binding",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d",
				Output:   g.Single(scb),
				Status:   http.StatusOK},
			Expected: scb,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceCredentialBindings.Get(context.Background(), "59ba6d78-6a21-4321-83a9-f7eacd88b08d")
			},
		},
		{
			Description: "Get service credential binding detail",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d/details",
				Output:   g.Single(scbd),
				Status:   http.StatusOK},
			Expected: scbd,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceCredentialBindings.GetDetails(context.Background(), "59ba6d78-6a21-4321-83a9-f7eacd88b08d")
			},
		},
		{
			Description: "Get service credential binding parameters",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d/parameters",
				Output:   g.Single(`{ "foo": "bar", "foz": "baz" }`),
				Status:   http.StatusOK},
			Expected: `{ "foo": "bar", "foz": "baz" }`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceCredentialBindings.GetParameters(context.Background(), "59ba6d78-6a21-4321-83a9-f7eacd88b08d")
			},
		},
		{
			Description: "Get service credential binding and app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource: scb,
					Apps:     []string{app},
				}),
				Status: http.StatusOK,
			},
			Expected:  scb,
			Expected2: app,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceCredentialBindings.GetIncludeApp(context.Background(), "59ba6d78-6a21-4321-83a9-f7eacd88b08d")
			},
		},
		{
			Description: "Get service credential binding and service instance",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource:         scb,
					ServiceInstances: []string{si},
				}),
				Status: http.StatusOK,
			},
			Expected:  scb,
			Expected2: si,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceCredentialBindings.GetIncludeServiceInstance(context.Background(), "59ba6d78-6a21-4321-83a9-f7eacd88b08d")
			},
		},
		{
			Description: "List all service credential bindings",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings",
				Output:   g.Paged([]string{scb}, []string{scb2}),
				Status:   http.StatusOK},
			Expected: g.Array(scb, scb2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceCredentialBindings.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List all service credential bindings and include apps",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources: []string{scb, scb2},
						Apps:      []string{app, app2},
					},
					testutil.PagedResult{
						Resources: []string{scb3, scb4},
						Apps:      []string{app3, app4},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(scb, scb2, scb3, scb4),
			Expected2: g.Array(app, app2, app3, app4),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceCredentialBindings.ListIncludeAppsAll(context.Background(), nil)
			},
		},
		{
			Description: "List all service credential bindings and include service instances",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_credential_bindings",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources:        []string{scb, scb2},
						ServiceInstances: []string{si},
					},
					testutil.PagedResult{
						Resources: []string{scb3, scb4},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(scb, scb2, scb3, scb4),
			Expected2: g.Array(si),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServiceCredentialBindings.ListIncludeServiceInstancesAll(context.Background(), nil)
			},
		},
		{
			Description: "Update service credential binding",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_credential_bindings/59ba6d78-6a21-4321-83a9-f7eacd88b08d",
				Output:   g.Single(scb),
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": {"foo": "bar"}, "annotations": {"baz": "qux"} }}`,
			},
			Expected: scb,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.ServiceCredentialBindingUpdate{
					Metadata: resource.NewMetadata().
						WithLabel("", "foo", "bar").
						WithAnnotation("", "baz", "qux"),
				}
				return c.ServiceCredentialBindings.Update(context.Background(), "59ba6d78-6a21-4321-83a9-f7eacd88b08d", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
