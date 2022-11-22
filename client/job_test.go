package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestJobs(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	job := g.Job("COMPLETE").JSON

	tests := []RouteTest{
		{
			Description: "Get job",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
				Output:   g.Single(job),
				Status:   http.StatusOK},
			Expected: job,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Jobs.Get(context.Background(), "c33a5caf-77e0-4d6e-b587-5555d339bc9a")
			},
		},
	}
	ExecuteTests(tests, t)
}
