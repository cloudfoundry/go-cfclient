package client

import (
	"net/http"
	"testing"
)

func TestSpaceFeatures(t *testing.T) {
	tests := []RouteTest{
		{
			Description: "Enable SSH for a space",
			Route: MockRoute{
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
				err := c.SpaceFeatures.EnableSSH("000d1e0c-218e-470b-b5db-84481b89fa92", true)
				return nil, err
			},
		},
	}
	executeTests(tests, t)
}
