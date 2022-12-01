package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestServicePlans(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(156)
	svcPlan := g.ServicePlan().JSON
	svcPlan2 := g.ServicePlan().JSON
	svcPlan3 := g.ServicePlan().JSON
	svcPlan4 := g.ServicePlan().JSON
	space := g.Space().JSON
	space2 := g.Space().JSON
	org := g.Organization().JSON
	svcOffering := g.ServiceOffering().JSON

	tests := []RouteTest{
		{
			Description: "Delete service plan",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.ServicePlans.Delete(context.Background(), "79aae221-b2a6-4aaa-a134-76f605af46c9")
				return nil, err
			},
		},
		{
			Description: "Get service plan",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9",
				Output:   g.Single(svcPlan),
				Status:   http.StatusOK},
			Expected: svcPlan,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServicePlans.Get(context.Background(), "79aae221-b2a6-4aaa-a134-76f605af46c9")
			},
		},
		{
			Description: "List all service plans",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_plans",
				Output:   g.Paged([]string{svcPlan}, []string{svcPlan2}),
				Status:   http.StatusOK},
			Expected: g.Array(svcPlan, svcPlan2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServicePlans.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List all service plans include service offerings",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_plans",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources:        []string{svcPlan, svcPlan2},
						ServiceOfferings: []string{svcOffering},
					},
					testutil.PagedResult{
						Resources: []string{svcPlan3, svcPlan4},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(svcPlan, svcPlan2, svcPlan3, svcPlan4),
			Expected2: g.Array(svcOffering),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServicePlans.ListIncludeServiceOfferingAll(context.Background(), nil)
			},
		},
		{
			Description: "List all service plans include spaces and organizations",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_plans",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources:     []string{svcPlan, svcPlan2},
						Spaces:        []string{space},
						Organizations: []string{org},
					},
					testutil.PagedResult{
						Resources: []string{svcPlan3, svcPlan4},
						Spaces:    []string{space2},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(svcPlan, svcPlan2, svcPlan3, svcPlan4),
			Expected2: g.Array(space, space2),
			Expected3: g.Array(org),
			Action3: func(c *Client, t *testing.T) (any, any, any, error) {
				return c.ServicePlans.ListIncludeSpacesAndOrganizationsAll(context.Background(), nil)
			},
		},
		{
			Description: "Update service plan",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9",
				Output:   g.Single(svcPlan),
				Status:   http.StatusOK,
				PostForm: `{
					"metadata": {
					  "labels": {"key": "value"},
					  "annotations": {"note": "detailed information"}
					}
				  }`,
			},
			Expected: svcPlan,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.ServicePlanUpdate{
					Metadata: resource.NewMetadata().
						WithLabel("", "key", "value").
						WithAnnotation("", "note", "detailed information"),
				}
				return c.ServicePlans.Update(context.Background(), "79aae221-b2a6-4aaa-a134-76f605af46c9", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
