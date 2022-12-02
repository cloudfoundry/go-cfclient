package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestServicePlanVisibilities(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(156)
	svcPlanVisibility := g.ServicePlanVisibility().JSON

	tests := []RouteTest{
		{
			Description: "Apply service plan visibility",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9/visibility",
				Output:   g.Single(svcPlanVisibility),
				Status:   http.StatusOK,
				PostForm: `{
					"type": "organization",
					"organizations": [
					  { "guid" : "0fc1ad4f-e1d7-4436-8e23-6b20f03c6482" }
					]
				  }`,
			},
			Expected: svcPlanVisibility,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServicePlanVisibilityUpdate(resource.ServicePlanVisibilityOrganization)
				r.Organizations = []resource.ServicePlanVisibilityRelation{
					{
						GUID: "0fc1ad4f-e1d7-4436-8e23-6b20f03c6482",
					},
				}
				return c.ServicePlansVisibility.Apply(context.Background(), "79aae221-b2a6-4aaa-a134-76f605af46c9", r)
			},
		},
		{
			Description: "Delete service plan visibility",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9/visibility/90a4d2ca-054b-4f15-9a44-cc94f845df9c",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.ServicePlansVisibility.Delete(context.Background(), "79aae221-b2a6-4aaa-a134-76f605af46c9", "90a4d2ca-054b-4f15-9a44-cc94f845df9c")
				return nil, err
			},
		},
		{
			Description: "Get service plan visibility",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9/visibility",
				Output:   g.Single(svcPlanVisibility),
				Status:   http.StatusOK,
			},
			Expected: svcPlanVisibility,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServicePlansVisibility.Get(context.Background(), "79aae221-b2a6-4aaa-a134-76f605af46c9")
			},
		},
		{
			Description: "Update service plan visibility",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9/visibility",
				Output:   g.Single(svcPlanVisibility),
				Status:   http.StatusOK,
				PostForm: `{
					"type": "organization",
					"organizations": [
					  { "guid" : "0fc1ad4f-e1d7-4436-8e23-6b20f03c6482" }
					]
				  }`,
			},
			Expected: svcPlanVisibility,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServicePlanVisibilityUpdate(resource.ServicePlanVisibilityOrganization)
				r.Organizations = []resource.ServicePlanVisibilityRelation{
					{
						GUID: "0fc1ad4f-e1d7-4436-8e23-6b20f03c6482",
					},
				}
				return c.ServicePlansVisibility.Update(context.Background(), "79aae221-b2a6-4aaa-a134-76f605af46c9", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
