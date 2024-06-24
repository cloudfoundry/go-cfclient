package client

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/cloudfoundry/go-cfclient/v3/testutil"
)

func TestBuildpacks(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1002)
	buildpack := g.Buildpack().JSON
	buildpack2 := g.Buildpack().JSON
	buildpack3 := g.Buildpack().JSON
	buildpack4 := g.Buildpack().JSON

	tests := []RouteTest{
		{
			Description: "Create buildpack",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/buildpacks",
				Output:   g.Single(buildpack),
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
				return c.Buildpacks.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete buildpack",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/buildpacks/6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Buildpacks.Delete(context.Background(), "6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb")
			},
		},
		{
			Description: "Get buildpack",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/buildpacks/6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb",
				Output:   g.Single(buildpack),
				Status:   http.StatusOK},
			Expected: buildpack,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Buildpacks.Get(context.Background(), "6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb")
			},
		},
		{
			Description: "List all buildpacks",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/buildpacks",
				Output:   g.Paged([]string{buildpack, buildpack2}, []string{buildpack3, buildpack4}),
				Status:   http.StatusOK},
			Expected: g.Array(buildpack, buildpack2, buildpack3, buildpack4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Buildpacks.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "Update buildpack",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/buildpacks/6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb",
				Output:   g.Single(buildpack),
				Status:   http.StatusOK,
				PostForm: `{ 
							"position": 1,
							"stack" : "cflinuxfs4"
							}`,
			},
			Expected: buildpack,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewBuildpackUpdate().
					WithPosition(1).
					WithStack("cflinuxfs4")
				return c.Buildpacks.Update(context.Background(), "6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb", r)
			},
		},
		{
			Description: "Upload buildpack",
			Route: testutil.MockRoute{
				Method:           "POST",
				Endpoint:         "/v3/buildpacks/6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb/upload",
				Output:           g.Single(buildpack),
				Status:           http.StatusOK,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected:  "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Expected2: buildpack,
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				zipFile := strings.NewReader("bp")
				return c.Buildpacks.Upload(context.Background(), "6f3c68d0-e119-4ca2-8ce4-83661ad6e0eb", "buildpack.zip", zipFile)
			},
		},
	}
	ExecuteTests(tests, t)
}
