package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestServicePlanVisibilities(t *testing.T) {
	g := test.NewObjectJSONGenerator(156)
	svcPlanVisibility := g.ServicePlanVisibility()

	tests := []RouteTest{
		{
			Description: "Apply service plan visibility",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9/visibility",
				Output:   []string{svcPlanVisibility},
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
				return c.ServicePlansVisibility.Apply("79aae221-b2a6-4aaa-a134-76f605af46c9", r)
			},
		},
		{
			Description: "Delete service plan visibility",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9/visibility/90a4d2ca-054b-4f15-9a44-cc94f845df9c",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.ServicePlansVisibility.Delete("79aae221-b2a6-4aaa-a134-76f605af46c9", "90a4d2ca-054b-4f15-9a44-cc94f845df9c")
				return nil, err
			},
		},
		{
			Description: "Get service plan visibility",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9/visibility",
				Output:   []string{svcPlanVisibility},
				Status:   http.StatusOK},
			Action: func(c *Client, t *testing.T) (any, error) {
				v, err := c.ServicePlansVisibility.Get("79aae221-b2a6-4aaa-a134-76f605af46c9")
				require.NoError(t, err)
				require.Equal(t, resource.ServicePlanVisibilityOrganization, v)
				return nil, nil
			},
		},
		{
			Description: "Update service plan visibility",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9/visibility",
				Output:   []string{svcPlanVisibility},
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
				return c.ServicePlansVisibility.Update("79aae221-b2a6-4aaa-a134-76f605af46c9", r)
			},
		},
	}
	executeTests(tests, t)
}
