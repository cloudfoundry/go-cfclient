package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestSpaceQuotas(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(15954)
	spaceQuota := g.SpaceQuota().JSON
	spaceQuota2 := g.SpaceQuota().JSON
	spaceQuota3 := g.SpaceQuota().JSON
	spaceQuota4 := g.SpaceQuota().JSON

	tests := []RouteTest{
		{
			Description: "Apply space quota to space",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/space_quotas/8a5955c0-d6fd-4f46-8e43-72a4dc35fb04/relationships/spaces",
				Output: []string{`{
					"data": [
					  { "guid": "ac79b04c-c9a2-488d-b830-3e5f26e600d1" },
					  { "guid": "284d7b6e-8447-40b3-8ab6-2b4926fca12d" }
					]
				  }`},
				Status:   http.StatusOK,
				PostForm: `{ "data": [{ "guid": "ac79b04c-c9a2-488d-b830-3e5f26e600d1" }, { "guid": "284d7b6e-8447-40b3-8ab6-2b4926fca12d" }] }`,
			},
			Expected: `["ac79b04c-c9a2-488d-b830-3e5f26e600d1", "284d7b6e-8447-40b3-8ab6-2b4926fca12d"]`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SpaceQuotas.Apply(context.Background(), "8a5955c0-d6fd-4f46-8e43-72a4dc35fb04", []string{
					"ac79b04c-c9a2-488d-b830-3e5f26e600d1", "284d7b6e-8447-40b3-8ab6-2b4926fca12d",
				})
			},
		},
		{
			Description: "Create space quota",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/space_quotas",
				Output:   g.Single(spaceQuota),
				Status:   http.StatusCreated,
				PostForm: `{
					"name": "my-space-quota",
					"relationships": {
					  "organization": {
						"data": {
						  "guid": "d6f5727f-c8a1-4f8e-93fb-440888b3bef1"
						}
					  }
					}
				  }`,
			},
			Expected: spaceQuota,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewSpaceQuotaCreate("my-space-quota", "d6f5727f-c8a1-4f8e-93fb-440888b3bef1")
				return c.SpaceQuotas.Create(context.Background(), r)
			},
		},
		{
			Description: "Get space quota",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/space_quotas/8a5955c0-d6fd-4f46-8e43-72a4dc35fb04",
				Output:   g.Single(spaceQuota),
				Status:   http.StatusOK,
			},
			Expected: spaceQuota,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SpaceQuotas.Get(context.Background(), "8a5955c0-d6fd-4f46-8e43-72a4dc35fb04")
			},
		},
		{
			Description: "Delete space quota",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/space_quotas/8a5955c0-d6fd-4f46-8e43-72a4dc35fb04",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SpaceQuotas.Delete(context.Background(), "8a5955c0-d6fd-4f46-8e43-72a4dc35fb04")
			},
		},
		{
			Description: "List all space quotas",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/space_quotas",
				Output:   g.Paged([]string{spaceQuota, spaceQuota2}, []string{spaceQuota3, spaceQuota4}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(spaceQuota, spaceQuota2, spaceQuota3, spaceQuota4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SpaceQuotas.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "Remove space quota",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/space_quotas/8a5955c0-d6fd-4f46-8e43-72a4dc35fb04/relationships/spaces/ac79b04c-c9a2-488d-b830-3e5f26e600d1",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.SpaceQuotas.Remove(context.Background(), "8a5955c0-d6fd-4f46-8e43-72a4dc35fb04", "ac79b04c-c9a2-488d-b830-3e5f26e600d1")
			},
		},
		{
			Description: "Update space quota",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/space_quotas/8a5955c0-d6fd-4f46-8e43-72a4dc35fb04",
				Output:   g.Single(spaceQuota),
				Status:   http.StatusOK,
				PostForm: `{
					"name": "don-quixote",
					"apps": {
					  "total_memory_in_mb": 5120,
					  "per_process_memory_in_mb": 1024,
					  "log_rate_limit_in_bytes_per_second": 1024,
					  "total_instances": 10,
					  "per_app_tasks": 5
					},
					"services": {
					  "paid_services_allowed": true,
					  "total_service_instances": 10,
					  "total_service_keys": 20
					},
					"routes": {
					  "total_routes": 8,
					  "total_reserved_ports": 4
					}
				  }`,
			},
			Expected: spaceQuota,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewSpaceQuotaUpdate().
					WithName("don-quixote").
					WithTotalMemoryInMB(5120).
					WithPerProcessMemoryInMB(1024).
					WithLogRateLimitInBytesPerSecond(1024).
					WithTotalInstances(10).
					WithPerAppTasks(5).
					WithPaidServicesAllowed(true).
					WithTotalServiceInstances(10).
					WithTotalServiceKeys(20).
					WithTotalRoutes(8).
					WithTotalReservedPorts(4)
				return c.SpaceQuotas.Update(context.Background(), "8a5955c0-d6fd-4f46-8e43-72a4dc35fb04", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
