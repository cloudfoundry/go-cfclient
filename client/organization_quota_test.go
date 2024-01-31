package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
)

func TestOrganizationQuotas(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(15)
	orgQuota := g.OrganizationQuota().JSON
	orgQuota2 := g.OrganizationQuota().JSON
	orgQuota3 := g.OrganizationQuota().JSON
	orgQuota4 := g.OrganizationQuota().JSON

	tests := []RouteTest{
		{
			Description: "Apply organization quota to org",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/organization_quotas/e3bff602-f3d4-4c63-a85a-d7155aa2f1ff/relationships/organizations",
				Output: []string{`{
					"data": [
					  { "guid": "5ab8881d-b5d7-40da-9228-b0253a7a1570" },
					  { "guid": "0610e399-1333-4d5f-b211-3a1135e0576e" }
					]
				  }`},
				Status:   http.StatusOK,
				PostForm: `{ "data": [{ "guid": "5ab8881d-b5d7-40da-9228-b0253a7a1570" }, { "guid": "0610e399-1333-4d5f-b211-3a1135e0576e" }] }`,
			},
			Expected: `["5ab8881d-b5d7-40da-9228-b0253a7a1570", "0610e399-1333-4d5f-b211-3a1135e0576e"]`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.OrganizationQuotas.Apply(context.Background(), "e3bff602-f3d4-4c63-a85a-d7155aa2f1ff", []string{
					"5ab8881d-b5d7-40da-9228-b0253a7a1570", "0610e399-1333-4d5f-b211-3a1135e0576e",
				})
			},
		},
		{
			Description: "Create organization quota",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/organization_quotas",
				Output:   g.Single(orgQuota),
				Status:   http.StatusCreated,
				PostForm: `{ 
					"name": "my-org-quota",
					"relationships": {
						"organizations": {
							"data": [
								{
									"guid": "e3bff602-f3d4-4c63-a85a-d7155aa2f1ff"
								}
							]
						}
					} 
				}`,
			},
			Expected: orgQuota,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewOrganizationQuotaCreate("my-org-quota")
				r.WithOrganizations("e3bff602-f3d4-4c63-a85a-d7155aa2f1ff")
				return c.OrganizationQuotas.Create(context.Background(), r)
			},
		},
		{
			Description: "Get organization quota",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organization_quotas/e3bff602-f3d4-4c63-a85a-d7155aa2f1ff",
				Output:   g.Single(orgQuota),
				Status:   http.StatusOK,
			},
			Expected: orgQuota,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.OrganizationQuotas.Get(context.Background(), "e3bff602-f3d4-4c63-a85a-d7155aa2f1ff")
			},
		},
		{
			Description: "Delete organization quota",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/organization_quotas/e3bff602-f3d4-4c63-a85a-d7155aa2f1ff",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.OrganizationQuotas.Delete(context.Background(), "e3bff602-f3d4-4c63-a85a-d7155aa2f1ff")
			},
		},
		{
			Description: "List all organization quotas",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organization_quotas",
				Output:   g.Paged([]string{orgQuota, orgQuota2}, []string{orgQuota3, orgQuota4}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(orgQuota, orgQuota2, orgQuota3, orgQuota4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.OrganizationQuotas.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "Update organization quota",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/organization_quotas/e3bff602-f3d4-4c63-a85a-d7155aa2f1ff",
				Output:   g.Single(orgQuota),
				Status:   http.StatusOK,
				PostForm: `{
					"name": "new_name",
					"apps": {
						"log_rate_limit_in_bytes_per_second": 1000,
						"per_app_tasks": 5,
						"per_process_memory_in_mb": 10,
						"total_instances": 15,
						"total_memory_in_mb": 100
					},
					"routes": {
						"total_reserved_ports": 35,
						"total_routes": 30
					},
					"services": {
						"paid_services_allowed": false,
						"total_service_instances": 20,
						"total_service_keys": 25
					}
				}`,
			},
			Expected: orgQuota,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewOrganizationQuotaUpdate().
					WithName("new_name").
					WithPerProcessMemoryInMB(10).
					WithAppsTotalMemoryInMB(100).
					WithTotalInstances(15).
					WithLogRateLimitInBytesPerSecond(1000).
					WithPerAppTasks(5).
					WithPaidServicesAllowed(false).
					WithTotalServiceInstances(20).
					WithTotalServiceKeys(25).
					WithTotalRoutes(30).
					WithTotalReservedPorts(35)
				return c.OrganizationQuotas.Update(context.Background(), "e3bff602-f3d4-4c63-a85a-d7155aa2f1ff", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
