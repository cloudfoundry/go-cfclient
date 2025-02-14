package client

import (
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"

	"github.com/cloudfoundry/go-cfclient/v3/testutil"
)

func TestJobs(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	job := g.Job("COMPLETE").JSON
	jobProcessing := g.Job("PROCESSING").JSON
	jobFailed := g.JobFailed().JSON
	pollingOpts := &PollingOptions{
		FailedState:   "FAILED",
		Timeout:       time.Second,
		CheckInterval: time.Nanosecond,
	}

	tests := []RouteTest{
		{
			Description: "Get job",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
				Output:   g.Single(job),
				Status:   http.StatusOK,
			},
			Expected: job,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Jobs.Get(context.Background(), "c33a5caf-77e0-4d6e-b587-5555d339bc9a")
			},
		},
		{
			Description: "Poll job that succeeds",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
				Output:   []string{jobProcessing, job},
				Statuses: []int{http.StatusOK, http.StatusOK},
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.Jobs.PollComplete(context.Background(), "c33a5caf-77e0-4d6e-b587-5555d339bc9a", pollingOpts)
				return nil, err
			},
		},
		{
			Description: "Poll job that fails",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/jobs/40e49716-44a3-4ae9-9926-d1a804acf70c",
				Output:   []string{jobProcessing, jobFailed},
				Statuses: []int{http.StatusOK, http.StatusOK},
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.Jobs.PollComplete(context.Background(), "40e49716-44a3-4ae9-9926-d1a804acf70c", pollingOpts)
				require.Error(t, err)
				require.ErrorContains(t, err, "received state FAILED while waiting for async process")
				require.ErrorContains(t, err, "cfclient error (CF-UnprocessableEntity|10008): something went wrong")
				require.ErrorContains(t, err, "cfclient error (UnknownError|10001): unexpected error occurred")
				return nil, nil
			},
		},
	}
	ExecuteTests(tests, t)
}
