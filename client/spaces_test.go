package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"net/http"
	"testing"
)

func TestSpaces(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	space := g.Space()
	space2 := g.Space()
	user := g.User()
	user2 := g.User()

	tests := []RouteTest{
		{
			Description: "Create space",
			Route: MockRoute{
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
			Route: MockRoute{
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
			Route: MockRoute{
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
			Description: "List all spaces",
			Route: MockRoute{
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
			Description: "List all space users",
			Route: MockRoute{
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
			Route: MockRoute{
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
	executeTests(tests, t)
}
