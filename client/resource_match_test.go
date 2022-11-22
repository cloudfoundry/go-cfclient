package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestResourceMatches(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	resourceMatch := g.ResourceMatch().JSON

	tests := []RouteTest{
		{
			Description: "Create a resource match",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/resource_matches",
				Output:   g.Single(resourceMatch),
				Status:   http.StatusOK,
				PostForm: `{
					"resources": [
					  {
						"checksum": { "value": "002d760bea1be268e27077412e11a320d0f164d3" },
						"size_in_bytes": 36,
						"path": "C:\\path\\to\\file",
						"mode": "645"
					  },
					  {
						"checksum": { "value": "a9993e364706816aba3e25717850c26c9cd0d89d" },
						"size_in_bytes": 1,
						"path": "path/to/file",
						"mode": "644"
					  }
					]
				  }`,
			},
			Expected: resourceMatch,
			Action: func(c *Client, t *testing.T) (any, error) {
				toMatch := &resource.ResourceMatches{
					Resources: []resource.ResourceMatch{
						{
							Path:        `C:\path\to\file`,
							Mode:        "645",
							SizeInBytes: 36,
							Checksum: resource.ResourceMatchChecksum{
								Value: "002d760bea1be268e27077412e11a320d0f164d3",
							},
						},
						{
							Path:        `path/to/file`,
							Mode:        "644",
							SizeInBytes: 1,
							Checksum: resource.ResourceMatchChecksum{
								Value: "a9993e364706816aba3e25717850c26c9cd0d89d",
							},
						},
					},
				}
				return c.ResourceMatches.Create(context.Background(), toMatch)
			},
		},
	}
	ExecuteTests(tests, t)
}
