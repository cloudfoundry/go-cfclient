package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestServiceOfferings(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(156)
	so := g.ServiceOffering().JSON
	so2 := g.ServiceOffering().JSON

	tests := []RouteTest{
		{
			Description: "Delete service offering",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_offerings/928a32d9-8101-4b86-85a4-96e06f833c2d",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.ServiceOfferings.Delete(context.Background(), "928a32d9-8101-4b86-85a4-96e06f833c2d")
				return nil, err
			},
		},
		{
			Description: "Get service offering",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_offerings/928a32d9-8101-4b86-85a4-96e06f833c2d",
				Output:   g.Single(so),
				Status:   http.StatusOK},
			Expected: so,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceOfferings.Get(context.Background(), "928a32d9-8101-4b86-85a4-96e06f833c2d")
			},
		},
		{
			Description: "List all service offerings",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_offerings",
				Output:   g.Paged([]string{so}, []string{so2}),
				Status:   http.StatusOK},
			Expected: g.Array(so, so2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceOfferings.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "Update service offering",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_offerings/928a32d9-8101-4b86-85a4-96e06f833c2d",
				Output:   g.Single(so),
				Status:   http.StatusOK,
				PostForm: `{
					"metadata": {
					  "labels": {"key": "value"},
					  "annotations": {"note": "detailed information"}
					}
				  }`,
			},
			Expected: so,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.ServiceOfferingUpdate{
					Metadata: resource.NewMetadata().
						WithLabel("", "key", "value").
						WithAnnotation("", "note", "detailed information"),
				}
				return c.ServiceOfferings.Update(context.Background(), "928a32d9-8101-4b86-85a4-96e06f833c2d", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
