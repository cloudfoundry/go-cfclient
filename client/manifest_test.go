package client

import (
	"github.com/cloudfoundry-community/go-cfclient/test"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestManifests(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	manifest := g.Manifest()

	tests := []RouteTest{
		{
			Description: "Generate app manifest",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/389f0d73-04ee-455b-b63c-513c7c78d5ff/manifest",
				Output:   []string{manifest},
				Status:   http.StatusOK},
			Action: func(c *Client, t *testing.T) (any, error) {
				actual, err := c.Manifests.Generate("389f0d73-04ee-455b-b63c-513c7c78d5ff")
				require.NoError(t, err)
				require.Equal(t, manifest, actual)
				return nil, nil
			},
		},
	}
	executeTests(tests, t)
}