package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestRevisions(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	revision := g.Revision().JSON
	revision2 := g.Revision().JSON
	revision3 := g.Revision().JSON
	revision4 := g.Revision().JSON
	appEnvVar := g.AppEnvVar().JSON

	tests := []RouteTest{
		{
			Description: "Get revision",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/revisions/5a49a370-92cd-4091-bb62-e0914460f7b2",
				Output:   g.Single(revision),
				Status:   http.StatusOK},
			Expected: revision,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Revisions.Get(context.Background(), "5a49a370-92cd-4091-bb62-e0914460f7b2")
			},
		},
		{
			Description: "Get revision environment variables",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/revisions/5a49a370-92cd-4091-bb62-e0914460f7b2/environment_variables",
				Output:   g.Single(appEnvVar),
				Status:   http.StatusOK,
			},
			Expected: `{ "RAILS_ENV": "production" }`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Revisions.GetEnvironmentVariables(context.Background(), "5a49a370-92cd-4091-bb62-e0914460f7b2")
			},
		},
		{
			Description: "List all app revisions",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/487d2a80-3769-4ad8-8ef5-a02c363d017b/revisions",
				Output:   g.Paged([]string{revision, revision2}, []string{revision3, revision4}),
				Status:   http.StatusOK},
			Expected: g.Array(revision, revision2, revision3, revision4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Revisions.ListForAppAll(context.Background(), "487d2a80-3769-4ad8-8ef5-a02c363d017b", nil)
			},
		},
		{
			Description: "List all deployed app revisions",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/487d2a80-3769-4ad8-8ef5-a02c363d017b/revisions/deployed",
				Output:   g.Paged([]string{revision, revision2}, []string{revision3, revision4}),
				Status:   http.StatusOK},
			Expected: g.Array(revision, revision2, revision3, revision4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Revisions.ListForAppDeployedAll(context.Background(), "487d2a80-3769-4ad8-8ef5-a02c363d017b", nil)
			},
		},
		{
			Description: "Update a revision",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/revisions/5a49a370-92cd-4091-bb62-e0914460f7b2",
				Output:   g.Single(revision),
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": { "key": "value" }, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: revision,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.RevisionUpdate{
					Metadata: resource.NewMetadata().
						WithLabel("", "key", "value").
						WithAnnotation("", "note", "detailed information"),
				}
				return c.Revisions.Update(context.Background(), "5a49a370-92cd-4091-bb62-e0914460f7b2", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
