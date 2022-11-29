package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
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
				PostForm: `{ "name": "my-org-quota" }`,
			},
			Expected: orgQuota,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewOrganizationQuotaCreate("my-org-quota")
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
				Method:   "DELETE",
				Endpoint: "/v3/organization_quotas/e3bff602-f3d4-4c63-a85a-d7155aa2f1ff",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.OrganizationQuotas.Delete(context.Background(), "e3bff602-f3d4-4c63-a85a-d7155aa2f1ff")
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
				PostForm: `{ "name": "new_name", "apps": { "per_app_tasks": 5 }}`,
			},
			Expected: orgQuota,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewOrganizationQuotaUpdate().WithName("new_name").WithPerAppTasks(5)
				return c.OrganizationQuotas.Update(context.Background(), "e3bff602-f3d4-4c63-a85a-d7155aa2f1ff", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
