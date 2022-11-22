package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestServiceUsages(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(161)
	serviceUsage := g.ServiceUsage().JSON
	serviceUsage2 := g.ServiceUsage().JSON
	serviceUsage3 := g.ServiceUsage().JSON

	tests := []RouteTest{
		{
			Description: "Get service usage event",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_usage_events/cb4fb5eb-9b72-4696-b7bc-666696dec1b3",
				Output:   g.Single(serviceUsage),
				Status:   http.StatusOK},
			Expected: serviceUsage,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceUsageEvents.Get(context.Background(), "cb4fb5eb-9b72-4696-b7bc-666696dec1b3")
			},
		},
		{
			Description: "List all service usage events",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_usage_events",
				Output:   g.Paged([]string{serviceUsage, serviceUsage2}, []string{serviceUsage3}),
				Status:   http.StatusOK},
			Expected: g.Array(serviceUsage, serviceUsage2, serviceUsage3),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceUsageEvents.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "Purge all service usage events",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_usage_events/actions/destructively_purge_all_and_reseed",
				Status:   http.StatusOK},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.ServiceUsageEvents.Purge(context.Background())
				return nil, err
			},
		},
	}
	ExecuteTests(tests, t)
}
