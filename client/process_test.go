package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestProcesses(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(78)
	process := g.Process().JSON
	process2 := g.Process().JSON
	process3 := g.Process().JSON
	process4 := g.Process().JSON
	processStats := g.ProcessStats().JSON

	tests := []RouteTest{
		{
			Description: "Get process",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/processes/ec4ff362-60c5-47a0-8246-2a134537c606",
				Output:   []string{process},
				Status:   http.StatusOK,
			},
			Expected: process,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Processes.Get("ec4ff362-60c5-47a0-8246-2a134537c606")
			},
		},
		{
			Description: "Get process stats",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/processes/ec4ff362-60c5-47a0-8246-2a134537c606/stats",
				Output:   []string{processStats},
				Status:   http.StatusOK,
			},
			Expected: processStats,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Processes.GetStats("ec4ff362-60c5-47a0-8246-2a134537c606")
			},
		},
		{
			Description: "List all processes",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/processes",
				Output:   g.Paged([]string{process, process2}, []string{process3, process4}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(process, process2, process3, process4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Processes.ListAll(nil)
			},
		},
		{
			Description: "List all processes for app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/2a550283-9245-493e-af36-5e4b8703f896/processes",
				Output:   g.Paged([]string{process, process2}, []string{process3, process4}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(process, process2, process3, process4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Processes.ListForAppAll("2a550283-9245-493e-af36-5e4b8703f896", nil)
			},
		},
		{
			Description: "Scale a process",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/processes/ec4ff362-60c5-47a0-8246-2a134537c606/actions/scale",
				Output:   []string{process},
				Status:   http.StatusOK,
				PostForm: `{
					"instances": 5,
					"memory_in_mb": 256,
					"disk_in_mb": 1024,
					"log_rate_limit_in_bytes_per_second": 1024
				  }`,
			},
			Expected: process,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewProcessScale().
					WithInstances(5).
					WithMemoryInMB(256).
					WithDiskInMB(1024).
					WithLogRateLimitInBytesPerSecond(1024)
				return c.Processes.Scale("ec4ff362-60c5-47a0-8246-2a134537c606", r)
			},
		},
		{
			Description: "Update a process",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/processes/ec4ff362-60c5-47a0-8246-2a134537c606",
				Output:   []string{process},
				Status:   http.StatusOK,
				PostForm: `{
					"command": "rackup",
					"health_check": {
						"type": "http",
						"data": {
							"timeout": 60
						}
					}
				}`,
			},
			Expected: process,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewProcessUpdate().
					WithCommand("rackup").
					WithHealthCheckType("http").
					WithHealthCheckTimeout(60)
				return c.Processes.Update("ec4ff362-60c5-47a0-8246-2a134537c606", r)
			},
		},
		{
			Description: "Terminate process",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/processes/ec4ff362-60c5-47a0-8246-2a134537c606/instances/0",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Processes.Terminate("ec4ff362-60c5-47a0-8246-2a134537c606", 0)
			},
		},
	}
	ExecuteTests(tests, t)
}
