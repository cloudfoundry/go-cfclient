package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestIsolationSegments(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	iso := g.IsolationSegment().JSON
	iso2 := g.IsolationSegment().JSON
	iso3 := g.IsolationSegment().JSON
	iso4 := g.IsolationSegment().JSON
	isoRelations := g.IsolationSegmentRelationships().JSON

	tests := []RouteTest{
		{
			Description: "Create isolation segment",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/isolation_segments",
				Output:   g.Single(iso),
				Status:   http.StatusCreated,
				PostForm: `{ "name": "my-iso" }`,
			},
			Expected: iso,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewIsolationSegmentCreate("my-iso")
				return c.IsolationSegments.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete iso",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.IsolationSegments.Delete(context.Background(), "a45d5da8-67dc-4523-b34b-ffa68b8d8821")
			},
		},
		{
			Description: "Entitle isolation segment for organization",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821/relationships/organizations",
				Output:   g.Single(isoRelations),
				Status:   http.StatusCreated,
				PostForm: `{ "data": [{ "guid":"5700e458-283d-4528-806f-c3509e038f05" }]}`,
			},
			Expected: isoRelations,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.IsolationSegments.EntitleOrganization(context.Background(), "a45d5da8-67dc-4523-b34b-ffa68b8d8821", "5700e458-283d-4528-806f-c3509e038f05")
			},
		},
		{
			Description: "Get iso",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821",
				Output:   g.Single(iso),
				Status:   http.StatusOK},
			Expected: iso,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.IsolationSegments.Get(context.Background(), "a45d5da8-67dc-4523-b34b-ffa68b8d8821")
			},
		},
		{
			Description: "List all isolation segments",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/isolation_segments",
				Output:   g.Paged([]string{iso, iso2}, []string{iso3, iso4}),
				Status:   http.StatusOK},
			Expected: g.Array(iso, iso2, iso3, iso4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.IsolationSegments.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List all isolation segment related organizations",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821/relationships/organizations",
				Output: []string{`{
				  "data": [
					{
					  "guid": "68d54d31-9b3a-463b-ba94-e8e4c32edbac"
					},
					{
					  "guid": "b19f6525-cbd3-4155-b156-dc0c2a431b4c"
					}
				  ],
				  "links": {
					"self": {
					  "href": "https://api.example.org/v3/isolation_segments/bdeg4371-cbd3-4155-b156-dc0c2a431b4c/relationships/organizations"
					},
					"related": {
					  "href": "https://api.example.org/v3/isolation_segments/bdeg4371-cbd3-4155-b156-dc0c2a431b4c/organizations"
					}
				  }
				}`},
				Status: http.StatusOK,
			},
			Expected: `["68d54d31-9b3a-463b-ba94-e8e4c32edbac", "b19f6525-cbd3-4155-b156-dc0c2a431b4c"]`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.IsolationSegments.ListOrganizationRelationships(context.Background(), "a45d5da8-67dc-4523-b34b-ffa68b8d8821")
			},
		},
		{
			Description: "List all isolation segment related spaces",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821/relationships/spaces",
				Output: []string{`{
				  "data": [
					{
					  "guid": "885735b5-aea4-4cf5-8e44-961af0e41920"
					},
					{
					  "guid": "d4c91047-7b29-4fda-b7f9-04033e5c9c9f"
					}
				  ],
				  "links": {
					"self": {
					  "href": "https://api.example.org/v3/isolation_segments/bdeg4371-cbd3-4155-b156-dc0c2a431b4c/relationships/organizations"
					},
					"related": {
					  "href": "https://api.example.org/v3/isolation_segments/bdeg4371-cbd3-4155-b156-dc0c2a431b4c/organizations"
					}
				  }
				}`},
				Status: http.StatusOK,
			},
			Expected: `["885735b5-aea4-4cf5-8e44-961af0e41920", "d4c91047-7b29-4fda-b7f9-04033e5c9c9f"]`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.IsolationSegments.ListSpaceRelationships(context.Background(), "a45d5da8-67dc-4523-b34b-ffa68b8d8821")
			},
		},
		{
			Description: "Revoke isolation segment for organization",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821/relationships/organizations/5700e458-283d-4528-806f-c3509e038f05",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.IsolationSegments.RevokeOrganization(context.Background(), "a45d5da8-67dc-4523-b34b-ffa68b8d8821", "5700e458-283d-4528-806f-c3509e038f05")
				return nil, err
			},
		},
		{
			Description: "Update iso",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821",
				Output:   g.Single(iso),
				Status:   http.StatusOK,
				PostForm: `{ "name": "new-name" }`,
			},
			Expected: iso,
			Action: func(c *Client, t *testing.T) (any, error) {
				name := "new-name"
				r := &resource.IsolationSegmentUpdate{
					Name: &name,
				}
				return c.IsolationSegments.Update(context.Background(), "a45d5da8-67dc-4523-b34b-ffa68b8d8821", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
