package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestUsers(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	user := g.User()
	user2 := g.User()

	tests := []RouteTest{
		{
			Description: "Create user",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/users",
				Output:   []string{user},
				Status:   http.StatusCreated,
				PostForm: `{ "guid": "3ebeaa8b-fd55-4724-a764-9f2231d8f7db" }`,
			},
			Expected: user,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.UserCreate{
					GUID: "3ebeaa8b-fd55-4724-a764-9f2231d8f7db",
				}
				return c.Users.Create(r)
			},
		},
		{
			Description: "Delete user",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/users/3ebeaa8b-fd55-4724-a764-9f2231d8f7db",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Users.Delete("3ebeaa8b-fd55-4724-a764-9f2231d8f7db")
			},
		},
		{
			Description: "Get user",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/users/3ebeaa8b-fd55-4724-a764-9f2231d8f7db",
				Output:   []string{user},
				Status:   http.StatusOK},
			Expected: user,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Users.Get("3ebeaa8b-fd55-4724-a764-9f2231d8f7db")
			},
		},
		{
			Description: "List all users",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/users",
				Output:   g.Paged([]string{user}, []string{user2}),
				Status:   http.StatusOK},
			Expected: g.Array(user, user2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Users.ListAll(nil)
			},
		},
		{
			Description: "Update user",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/users/3ebeaa8b-fd55-4724-a764-9f2231d8f7db",
				Output:   []string{user},
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": { "key": "value" }, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: user,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.UserUpdate{
					Metadata: &resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.Users.Update("3ebeaa8b-fd55-4724-a764-9f2231d8f7db", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
