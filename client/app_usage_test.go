package client

import (
	"github.com/cloudfoundry-community/go-cfclient/test"
	"net/http"
	"testing"
)

func TestAppUsages(t *testing.T) {
	g := test.NewObjectJSONGenerator(161)
	appUsage := g.AppUsage()
	appUsage2 := g.AppUsage()
	appUsage3 := g.AppUsage()

	tests := []RouteTest{
		{
			Description: "Get app usage event",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/app_usage_events/af846b67-e0c4-44eb-bfa8-ff30e902d710",
				Output:   []string{appUsage},
				Status:   http.StatusOK},
			Expected: appUsage,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AppUsageEvents.Get("af846b67-e0c4-44eb-bfa8-ff30e902d710")
			},
		},
		{
			Description: "List all app usage events",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/app_usage_events",
				Output:   g.Paged([]string{appUsage, appUsage2}, []string{appUsage3}),
				Status:   http.StatusOK},
			Expected: g.Array(appUsage, appUsage2, appUsage3),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AppUsageEvents.ListAll(nil)
			},
		},
		{
			Description: "Purge all app usage events",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/app_usage_events/actions/destructively_purge_all_and_reseed",
				Status:   http.StatusOK},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.AppUsageEvents.Purge()
				return nil, err
			},
		},
	}
	executeTests(tests, t)
}
