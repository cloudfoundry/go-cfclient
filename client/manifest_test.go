package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestManifests(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	manifest := g.Manifest().JSON

	tests := []RouteTest{
		{
			Description: "Generate app manifest",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/389f0d73-04ee-455b-b63c-513c7c78d5ff/manifest",
				Output:   g.Single(manifest),
				Status:   http.StatusOK},
			Action: func(c *Client, t *testing.T) (any, error) {
				actual, err := c.Manifests.Generate(context.Background(), "389f0d73-04ee-455b-b63c-513c7c78d5ff")
				require.NoError(t, err)
				require.Equal(t, manifest, actual)
				return nil, nil
			},
		},
	}
	ExecuteTests(tests, t)
}
