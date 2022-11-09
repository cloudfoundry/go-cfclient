package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"net/http"
	"testing"
)

func TestServiceOfferings(t *testing.T) {
	g := test.NewObjectJSONGenerator(156)
	so := g.ServiceOffering()
	so2 := g.ServiceOffering()

	tests := []RouteTest{
		{
			Description: "Delete service offering",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_offerings/928a32d9-8101-4b86-85a4-96e06f833c2d",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.ServiceOfferings.Delete("928a32d9-8101-4b86-85a4-96e06f833c2d")
				return nil, err
			},
		},
		{
			Description: "Get service offering",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_offerings/928a32d9-8101-4b86-85a4-96e06f833c2d",
				Output:   []string{so},
				Status:   http.StatusOK},
			Expected: so,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceOfferings.Get("928a32d9-8101-4b86-85a4-96e06f833c2d")
			},
		},
		{
			Description: "List all service offerings",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_offerings",
				Output:   g.Paged([]string{so}, []string{so2}),
				Status:   http.StatusOK},
			Expected: g.Array(so, so2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceOfferings.ListAll(nil)
			},
		},
		{
			Description: "Update service offering",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_offerings/928a32d9-8101-4b86-85a4-96e06f833c2d",
				Output:   []string{so},
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
					Metadata: resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.ServiceOfferings.Update("928a32d9-8101-4b86-85a4-96e06f833c2d", r)
			},
		},
	}
	executeTests(tests, t)
}
