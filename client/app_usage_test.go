package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/testutil"
)

func TestAppUsages(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(161)
	appUsage := g.AppUsage().JSON
	appUsage2 := g.AppUsage().JSON
	appUsage3 := g.AppUsage().JSON

	tests := []RouteTest{
		{
			Description: "Get app usage event",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/app_usage_events/af846b67-e0c4-44eb-bfa8-ff30e902d710",
				Output:   g.Single(appUsage),
				Status:   http.StatusOK},
			Expected: appUsage,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AppUsageEvents.Get(context.Background(), "af846b67-e0c4-44eb-bfa8-ff30e902d710")
			},
		},
		{
			Description: "List all app usage events",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/app_usage_events",
				Output:   g.Paged([]string{appUsage, appUsage2}, []string{appUsage3}),
				Status:   http.StatusOK},
			Expected: g.Array(appUsage, appUsage2, appUsage3),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AppUsageEvents.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "Purge all app usage events",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/app_usage_events/actions/destructively_purge_all_and_reseed",
				Status:   http.StatusOK},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.AppUsageEvents.Purge(context.Background())
				return nil, err
			},
		},
	}
	ExecuteTests(tests, t)
}
