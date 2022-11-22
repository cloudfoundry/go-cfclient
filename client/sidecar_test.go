package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestSidecars(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(195)
	sidecar := g.Sidecar().JSON
	sidecar2 := g.Sidecar().JSON
	sidecar3 := g.Sidecar().JSON
	sidecar4 := g.Sidecar().JSON

	tests := []RouteTest{
		{
			Description: "Create a sidecar",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/apps/631b46a1-c3b6-4599-9659-72c9fd54817f/sidecars",
				Output:   g.Single(sidecar),
				Status:   http.StatusCreated,
				PostForm: `{
					"name": "auth-sidecar",
					"command": "bundle exec rackup",
					"process_types": ["web", "worker"],
					"memory_in_mb": 300
				  }`,
			},
			Expected: sidecar,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewSidecarCreate("auth-sidecar", "bundle exec rackup", []string{"web", "worker"}).
					WithMemoryInMB(300)
				return c.Sidecars.Create(context.Background(), "631b46a1-c3b6-4599-9659-72c9fd54817f", r)
			},
		},
		{
			Description: "Delete sidecar",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/sidecars/319ac7e8-e34a-4b6f-89da-1753ad3ece93",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Sidecars.Delete(context.Background(), "319ac7e8-e34a-4b6f-89da-1753ad3ece93")
			},
		},
		{
			Description: "Get sidecar",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/sidecars/319ac7e8-e34a-4b6f-89da-1753ad3ece93",
				Output:   g.Single(sidecar),
				Status:   http.StatusOK},
			Expected: sidecar,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Sidecars.Get(context.Background(), "319ac7e8-e34a-4b6f-89da-1753ad3ece93")
			},
		},
		{
			Description: "List all sidecars for app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/631b46a1-c3b6-4599-9659-72c9fd54817f/sidecars",
				Output:   g.Paged([]string{sidecar, sidecar2}, []string{sidecar3, sidecar4}),
				Status:   http.StatusOK},
			Expected: g.Array(sidecar, sidecar2, sidecar3, sidecar4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Sidecars.ListForAppAll(context.Background(), "631b46a1-c3b6-4599-9659-72c9fd54817f", nil)
			},
		},
		{
			Description: "List all sidecars for process",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/processes/0d2da177-c801-42a0-a6ca-ee4b10334954/sidecars",
				Output:   g.Paged([]string{sidecar, sidecar2}, []string{sidecar3, sidecar4}),
				Status:   http.StatusOK},
			Expected: g.Array(sidecar, sidecar2, sidecar3, sidecar4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Sidecars.ListForProcessAll(context.Background(), "0d2da177-c801-42a0-a6ca-ee4b10334954", nil)
			},
		},
	}
	ExecuteTests(tests, t)
}
