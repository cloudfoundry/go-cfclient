package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/testutil"
)

func TestSpaceFeatures(t *testing.T) {
	tests := []RouteTest{
		{
			Description: "Enable SSH for a space",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/spaces/000d1e0c-218e-470b-b5db-84481b89fa92/features/ssh",
				Output: []string{`{
					  "name": "ssh",
					  "description": "Enable SSHing into apps in the space.",
					  "enabled": true
					}`},
				Status:   http.StatusOK,
				PostForm: `{ "enabled": true }`,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.SpaceFeatures.EnableSSH(context.Background(), "000d1e0c-218e-470b-b5db-84481b89fa92", true)
				return nil, err
			},
		},
		{
			Description: "Is SSH enabled for a space",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/spaces/000d1e0c-218e-470b-b5db-84481b89fa92/features/ssh",
				Output: []string{`{
					  "name": "ssh",
					  "description": "Enable SSHing into apps in the space.",
					  "enabled": true
					}`},
				Status: http.StatusOK,
			},
			Expected: "true",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.SpaceFeatures.IsSSHEnabled(context.Background(), "000d1e0c-218e-470b-b5db-84481b89fa92")
			},
		},
	}
	ExecuteTests(tests, t)
}
