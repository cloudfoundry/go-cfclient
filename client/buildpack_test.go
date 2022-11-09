package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"net/http"
	"testing"
)

func TestBuildpacks(t *testing.T) {
	g := test.NewObjectJSONGenerator(1002)
	buildpack := g.Buildpack()
	buildpack2 := g.Buildpack()
	buildpack3 := g.Buildpack()
	buildpack4 := g.Buildpack()

	tests := []RouteTest{
		{
			Description: "Create buildpack",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/buildpacks",
				Output:   []string{buildpack},
				Status:   http.StatusCreated,
				PostForm: `{
					"name": "ruby_buildpack",
					"position": 42,
					"enabled": true,
					"locked": false,
					"stack": "cflinuxfs3"
				  }`,
			},
			Expected: buildpack,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewBuildpackCreate("ruby_buildpack").
					WithEnabled(true).
					WithPosition(42).
					WithLocked(false).
					WithStack("cflinuxfs3")
				return c.Buildpacks.Create(r)
			},
		},
		{
			Description: "Delete buildpack",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/buildpacks/6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Buildpacks.Delete("6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb")
			},
		},
		{
			Description: "Get buildpack",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/buildpacks/6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb",
				Output:   []string{buildpack},
				Status:   http.StatusOK},
			Expected: buildpack,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Buildpacks.Get("6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb")
			},
		},
		{
			Description: "List all buildpacks",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/buildpacks",
				Output:   g.Paged([]string{buildpack, buildpack2}, []string{buildpack3, buildpack4}),
				Status:   http.StatusOK},
			Expected: g.Array(buildpack, buildpack2, buildpack3, buildpack4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Buildpacks.ListAll(nil)
			},
		},
		{
			Description: "Update buildpack",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/buildpacks/6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb",
				Output:   []string{buildpack},
				Status:   http.StatusOK,
				PostForm: `{ "position": 1 }`,
			},
			Expected: buildpack,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewBuildpackUpdate().WithPosition(1)
				return c.Buildpacks.Update("6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb", r)
			},
		},
	}
	executeTests(tests, t)
}
