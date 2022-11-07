package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"net/http"
	"testing"
)

func TestRevisions(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	revision := g.Revision()
	revision2 := g.Revision()
	revision3 := g.Revision()
	revision4 := g.Revision()
	appEnvVar := g.AppEnvVar()

	tests := []RouteTest{
		{
			Description: "Get revision",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/revisions/5a49a370-92cd-4091-bb62-e0914460f7b2",
				Output:   []string{revision},
				Status:   http.StatusOK},
			Expected: revision,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Revisions.Get("5a49a370-92cd-4091-bb62-e0914460f7b2")
			},
		},
		{
			Description: "Get revision environment variables",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/revisions/5a49a370-92cd-4091-bb62-e0914460f7b2/environment_variables",
				Output:   []string{appEnvVar},
				Status:   http.StatusOK,
			},
			Expected: `{ "RAILS_ENV": "production" }`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Revisions.GetEnvironmentVariables("5a49a370-92cd-4091-bb62-e0914460f7b2")
			},
		},
		{
			Description: "List all app revisions",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/487d2a80-3769-4ad8-8ef5-a02c363d017b/revisions",
				Output:   g.Paged([]string{revision, revision2}, []string{revision3, revision4}),
				Status:   http.StatusOK},
			Expected: g.Array(revision, revision2, revision3, revision4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Revisions.ListAll("487d2a80-3769-4ad8-8ef5-a02c363d017b", nil)
			},
		},
		{
			Description: "List all deployed app revisions",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/487d2a80-3769-4ad8-8ef5-a02c363d017b/revisions/deployed",
				Output:   g.Paged([]string{revision, revision2}, []string{revision3, revision4}),
				Status:   http.StatusOK},
			Expected: g.Array(revision, revision2, revision3, revision4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Revisions.ListDeployedAll("487d2a80-3769-4ad8-8ef5-a02c363d017b", nil)
			},
		},
		{
			Description: "Update a revision",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/revisions/5a49a370-92cd-4091-bb62-e0914460f7b2",
				Output:   []string{revision},
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": { "key": "value" }, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: revision,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.RevisionUpdate{
					Metadata: &resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.Revisions.Update("5a49a370-92cd-4091-bb62-e0914460f7b2", r)
			},
		},
	}
	executeTests(tests, t)
}