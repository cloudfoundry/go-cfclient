package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"net/http"
	"testing"
)

func TestJobs(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	job := g.Job()

	tests := []RouteTest{
		{
			Description: "Get job",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
				Output:   []string{job},
				Status:   http.StatusOK},
			Expected: job,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Jobs.Get("c33a5caf-77e0-4d6e-b587-5555d339bc9a")
			},
		},
	}
	executeTests(tests, t)
}
