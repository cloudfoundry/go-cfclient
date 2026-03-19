package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/testutil"
)

func TestInfo(t *testing.T) {
	g := testutil.NewObjectJSONGenerator()
	info := g.Info().JSON
	usageSummary := g.InfoUsageSummary().JSON

	tests := []RouteTest{
		{
			Description: "Get platform info",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/info",
				Output:   g.Single(info),
				Status:   http.StatusOK},
			Expected: info,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Info.Get(context.Background())
			},
		},
		{
			Description: "Get platform usage summary",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/info/usage_summary",
				Output:   g.Single(usageSummary),
				Status:   http.StatusOK},
			Expected: usageSummary,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Info.GetUsageSummary(context.Background())
			},
		},
	}
	ExecuteTests(tests, t)
}
