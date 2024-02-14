package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
)

func TestServiceInstances(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(156)
	si := g.ServiceInstance().JSON
	si2 := g.ServiceInstance().JSON
	siSharedSummary := g.ServiceInstanceUsageSummary().JSON
	siSpaceRelationships := g.ServiceInstanceSpaceRelationships().JSON

	tests := []RouteTest{
		{
			Description: "Create managed service instance",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_instances",
				Output:   g.Single(si),
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
				RedirectLocation: "https://api.example.org/v3/jobs/af5c57f6-8769-41fa-a499-2c84ed896788",
			},
			Expected: "af5c57f6-8769-41fa-a499-2c84ed896788",
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServiceInstanceCreateManaged("my_service_instance",
					"7304bc3c-7010-11ea-8840-48bf6bec2d78", "e0e4417c-74ee-11ea-a604-48bf6bec2d78")
				r.Tags = []string{"foo", "bar", "baz"}
				return c.ServiceInstances.CreateManaged(context.Background(), r)
			},
		},
		{
			Description: "Create managed service instance with parameters",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_instances",
				Output:   g.Single(si),
				Status:   http.StatusCreated,
				PostForm: `{
					"type": "managed",
					"name": "my_service_instance",
					"tags": ["foo", "bar", "baz"],
					"parameters": {
						"foo": "bar",
						"baz": "qux"
					},
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
				RedirectLocation: "https://api.example.org/v3/jobs/af5c57f6-8769-41fa-a499-2c84ed896788",
			},
			Expected: "af5c57f6-8769-41fa-a499-2c84ed896788",
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServiceInstanceCreateManaged("my_service_instance",
					"7304bc3c-7010-11ea-8840-48bf6bec2d78", "e0e4417c-74ee-11ea-a604-48bf6bec2d78").
					WithTags([]string{"foo", "bar", "baz"}).
					WithParameters(json.RawMessage(`{"foo": "bar", "baz": "qux"}`))

				return c.ServiceInstances.CreateManaged(context.Background(), r)
			},
		},
		{
			Description: "Create user provided service instance",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_instances",
				Output:   g.Single(si),
				Status:   http.StatusCreated,
				PostForm: `{
					"type": "user-provided",
					"name": "my_service_instance",
					"tags": ["foo", "bar", "baz"],
					"relationships": {
						"space": {
							"data": {
								"guid": "7304bc3c-7010-11ea-8840-48bf6bec2d78"
							}
						}
					}
				}`,
			},
			Expected: si,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServiceInstanceCreateUserProvided("my_service_instance",
					"7304bc3c-7010-11ea-8840-48bf6bec2d78")
				r.Tags = []string{"foo", "bar", "baz"}
				return c.ServiceInstances.CreateUserProvided(context.Background(), r)
			},
		},
		{
			Description: "Create user provided service instance with credentials",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_instances",
				Output:   g.Single(si),
				Status:   http.StatusCreated,
				PostForm: `{
					"type": "user-provided",
					"name": "my_service_instance",
					"tags": ["foo", "bar", "baz"],
					"credentials": {
						"foo": "bar",
						"baz": "qux"
					},
					"relationships": {
						"space": {
							"data": {
								"guid": "7304bc3c-7010-11ea-8840-48bf6bec2d78"
							}
						}
					}
				}`,
			},
			Expected: si,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServiceInstanceCreateUserProvided("my_service_instance",
					"7304bc3c-7010-11ea-8840-48bf6bec2d78").
					WithCredentials(json.RawMessage(`{"foo": "bar", "baz": "qux"}`)).
					WithTags([]string{"foo", "bar", "baz"})
				return c.ServiceInstances.CreateUserProvided(context.Background(), r)
			},
		},
		{
			Description: "Delete service instance",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/v3/jobs/af5c57f6-8769-41fa-a499-2c84ed896788",
			},
			Expected: "af5c57f6-8769-41fa-a499-2c84ed896788",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceInstances.Delete(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8")
			},
		},
		{
			Description: "Get service instance",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8",
				Output:   g.Single(si),
				Status:   http.StatusOK},
			Expected: si,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceInstances.Get(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8")
			},
		},
		{
			Description: "Get service instance shared space relationships",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8/relationships/shared_spaces",
				Output:   g.Single(siSpaceRelationships),
				Status:   http.StatusOK},
			Expected: siSpaceRelationships,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceInstances.GetSharedSpaceRelationships(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8")
			},
		},
		{
			Description: "Get service instance shared space usage summary",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8/relationships/shared_spaces/usage_summary",
				Output:   g.Single(siSharedSummary),
				Status:   http.StatusOK},
			Expected: siSharedSummary,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceInstances.GetSharedSpaceUsageSummary(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8")
			},
		},
		{
			Description: "Get service instance user permissions",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8/permissions",
				Output:   g.Single(`{ "read": true, "manage": false }`),
				Status:   http.StatusOK},
			Expected: `{ "read": true, "manage": false }`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceInstances.GetUserPermissions(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8")
			},
		},
		{
			Description: "Get managed service instance parameters",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8/parameters",
				Output:   g.Single(`{ "key_1": "value_1", "key_2": "value_2" }`),
				Status:   http.StatusOK},
			Action: func(c *Client, t *testing.T) (any, error) {
				credentials, err := c.ServiceInstances.GetManagedParameters(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8")
				require.NoError(t, err)
				b, err := credentials.MarshalJSON()
				require.NoError(t, err)
				require.Equal(t, `{ "key_1": "value_1", "key_2": "value_2" }`, string(b))
				return nil, nil
			},
		},
		{
			Description: "Get user provided service instance credentials",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8/credentials",
				Output:   g.Single(`{ "username": "my-username", "password": "super-secret", "other": "credential" }`),
				Status:   http.StatusOK},
			Action: func(c *Client, t *testing.T) (any, error) {
				credentials, err := c.ServiceInstances.GetUserProvidedCredentials(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8")
				require.NoError(t, err)
				b, err := credentials.MarshalJSON()
				require.NoError(t, err)
				require.Equal(t, `{ "username": "my-username", "password": "super-secret", "other": "credential" }`, string(b))
				return nil, nil
			},
		},
		{
			Description: "List all service instances",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_instances",
				Output:   g.Paged([]string{si}, []string{si2}),
				Status:   http.StatusOK},
			Expected: g.Array(si, si2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceInstances.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "Update user provided service instance",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8",
				Output:   g.Single(si),
				Status:   http.StatusOK,
				PostForm: `{
					"name": "my_service_instance",
					"credentials": {
					  "foo": "bar",
					  "baz": "qux"
					},
					"tags": ["foo", "bar", "baz"],
					"syslog_drain_url": "https://syslog.com/drain",
					"route_service_url": "https://route.com/service"
				  }`,
			},
			Expected: si,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServiceInstanceUserProvidedUpdate().
					WithName("my_service_instance").
					WithCredentials(json.RawMessage(`{"foo": "bar", "baz": "qux"}`)).
					WithTags([]string{"foo", "bar", "baz"}).
					WithSyslogDrainURL("https://syslog.com/drain").
					WithRouteServiceURL("https://route.com/service")
				return c.ServiceInstances.UpdateUserProvided(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8", r)
			},
		},
		{
			Description: "Update managed service instance async",
			Route: testutil.MockRoute{
				Method:           "PATCH",
				Endpoint:         "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/v3/jobs/af5c57f6-8769-41fa-a499-2c84ed896788",
				PostForm: `{
					"name": "my_service_instance",
					"parameters": {
					  "foo": "bar",
					  "baz": "qux"
					},
					"tags": ["foo", "bar", "baz"],
					"relationships": {
					  "service_plan": {
						"data": {
						  "guid": "f2b6ba9c-a4d2-11ea-8ae6-48bf6bec2d78"
						}
					  }
					}
				  }`,
			},
			Expected: "af5c57f6-8769-41fa-a499-2c84ed896788",
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				r := resource.NewServiceInstanceManagedUpdate().
					WithName("my_service_instance").
					WithParameters(json.RawMessage(`{"foo": "bar", "baz": "qux"}`)).
					WithTags([]string{"foo", "bar", "baz"}).
					WithServicePlan("f2b6ba9c-a4d2-11ea-8ae6-48bf6bec2d78")
				return c.ServiceInstances.UpdateManaged(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8", r)
			},
		},
		{
			Description: "Update managed service instance sync",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8",
				Output:   g.Single(si),
				Status:   http.StatusOK,
				PostForm: `{ "name": "my_service_instance" }`,
			},
			Expected:  "",
			Expected2: si,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				r := resource.NewServiceInstanceManagedUpdate().
					WithName("my_service_instance")
				return c.ServiceInstances.UpdateManaged(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8", r)
			},
		},
		{
			Description: "Share service instance with space",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8/relationships/shared_spaces",
				Output:   g.Single(siSpaceRelationships),
				Status:   http.StatusOK,
				PostForm: `{ "data": [{ "guid":"000d1e0c-218e-470b-b5db-84481b89fa92" }]}`,
			},
			Expected: siSpaceRelationships,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceInstances.ShareWithSpace(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8", "000d1e0c-218e-470b-b5db-84481b89fa92")
			},
		},
		{
			Description: "Un-Share service instance with spaces",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_instances/62a3c0fe-5751-4f8f-97c4-28de85962ef8/relationships/shared_spaces/000d1e0c-218e-470b-b5db-84481b89fa92",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.ServiceInstances.UnShareWithSpaces(context.Background(), "62a3c0fe-5751-4f8f-97c4-28de85962ef8", []string{"000d1e0c-218e-470b-b5db-84481b89fa92"})
			},
		},
	}
	ExecuteTests(tests, t)
}
