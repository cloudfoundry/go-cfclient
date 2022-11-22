package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestAdmin(t *testing.T) {
	tests := []RouteTest{
		{
			Description: "Clear buildpack cache",
			Route: testutil.MockRoute{
				Method:           "POST",
				Endpoint:         "/v3/admin/actions/clear_buildpack_cache",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Admin.ClearBuildpackCache(context.Background())
			},
		},
	}
	ExecuteTests(tests, t)
}
