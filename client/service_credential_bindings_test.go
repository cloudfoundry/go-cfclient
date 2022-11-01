package client

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestServiceCredentialBindings(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	scb := g.ServiceCredentialBinding()
	scb2 := g.ServiceCredentialBinding()

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
