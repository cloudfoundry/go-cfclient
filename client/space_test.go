package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestSpaces(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	space := g.Space()
	space2 := g.Space()
	space3 := g.Space()
	space4 := g.Space()
	user := g.User()
	user2 := g.User()
	org := g.Organization()
	org2 := g.Organization()

	tests := []RouteTest{
		{
			Description: "Assign space iso segment",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/spaces/000d1e0c-218e-470b-b5db-84481b89fa92/relationships/isolation_segment",
				Output:   []string{`{ "data": { "guid": "443a1ea0-2403-4f0f-8c74-023a320bd1f2" }}`},
				Status:   http.StatusOK,
				PostForm: `{ "data": { "guid": "443a1ea0-2403-4f0f-8c74-023a320bd1f2" }}`,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.Spaces.AssignIsoSegment("000d1e0c-218e-470b-b5db-84481b89fa92", "443a1ea0-2403-4f0f-8c74-023a320bd1f2")
				return nil, err
			},
		},
		{
			Description: "Create space",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/spaces",
				Output:   []string{space},
				Status:   http.StatusCreated,
				PostForm: `{
					"name": "my-space",
					"relationships": {
						"organization": {
							"data": {
								"guid": "70c727ac-eef9-4e2a-aac3-975d5a0a0f15"
							}
						}
					}
				}`,
			},
			Expected: space,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewSpaceCreate("my-space", "70c727ac-eef9-4e2a-aac3-975d5a0a0f15")
				return c.Spaces.Create(r)
			},
		},
		{
			Description: "Delete space",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/spaces/000d1e0c-218e-470b-b5db-84481b89fa92",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Spaces.Delete("000d1e0c-218e-470b-b5db-84481b89fa92")
			},
		},
		{
			Description: "Get space",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/spaces/000d1e0c-218e-470b-b5db-84481b89fa92",
				Output:   []string{space},
				Status:   http.StatusOK},
			Expected: space,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Spaces.Get("000d1e0c-218e-470b-b5db-84481b89fa92")
			},
		},
		{
			Description: "Get assigned isolation segment",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/spaces/000d1e0c-218e-470b-b5db-84481b89fa92/relationships/isolation_segment",
				Output: []string{`{
					  "data": {
						"guid": "e4c91047-3b29-4fda-b7f9-04033e5a9c9f"
					  },
					  "links": {
						"self": {
						  "href": "https://api.example.org/v3/spaces/885735b5-aea4-4cf5-8e44-961af0e41920/relationships/isolation_segment"
						},
						"related": {
						  "href": "https://api.example.org/v3/isolation_segments/e4c91047-3b29-4fda-b7f9-04033e5a9c9f"
						}
					  }
					}`},
				Status: http.StatusOK},
			Expected: "e4c91047-3b29-4fda-b7f9-04033e5a9c9f",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Spaces.GetAssignedIsoSegment("000d1e0c-218e-470b-b5db-84481b89fa92")
			},
		},
		{
			Description: "Get space and org",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/spaces/000d1e0c-218e-470b-b5db-84481b89fa92",
				Output: g.ResourceWithInclude(testutil.ResourceResult{
					Resource:      space,
					Organizations: []string{org},
				}),
				Status: http.StatusOK,
			},
			Expected:  space,
			Expected2: org,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Spaces.GetIncludeOrg("000d1e0c-218e-470b-b5db-84481b89fa92")
			},
		},
		{
			Description: "List all spaces",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/spaces",
				Output:   g.Paged([]string{space}, []string{space2}),
				Status:   http.StatusOK},
			Expected: g.Array(space, space2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Spaces.ListAll(nil)
			},
		},
		{
			Description: "List all spaces and include parent orgs",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/spaces",
				Output: g.PagedWithInclude(
					testutil.PagedResult{
						Resources:     []string{space, space2},
						Organizations: []string{org},
					},
					testutil.PagedResult{
						Resources:     []string{space3, space4},
						Organizations: []string{org2},
					}),
				Status: http.StatusOK,
			},
			Expected:  g.Array(space, space2, space3, space4),
			Expected2: g.Array(org, org2),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.Spaces.ListIncludeOrgsAll(nil)
			},
		},
		{
			Description: "List all space users",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/spaces/000d1e0c-218e-470b-b5db-84481b89fa92/users",
				Output:   g.Paged([]string{user}, []string{user2}),
				Status:   http.StatusOK},
			Expected: g.Array(user, user2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Spaces.ListUsersAll("000d1e0c-218e-470b-b5db-84481b89fa92", nil)
			},
		},
		{
			Description: "Update space",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/spaces/000d1e0c-218e-470b-b5db-84481b89fa92",
				Output:   []string{space},
				Status:   http.StatusOK,
				PostForm: `{ "name": "new-space-name" }`,
			},
			Expected: space,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.SpaceUpdate{
					Name: "new-space-name",
				}
				return c.Spaces.Update("000d1e0c-218e-470b-b5db-84481b89fa92", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
