package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestStacks(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(165123)
	stack := g.Stack().JSON
	stack2 := g.Stack().JSON
	app := g.Application().JSON
	app2 := g.Application().JSON

	tests := []RouteTest{
		{
			Description: "Create stack",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/stacks",
				Output:   g.Single(stack),
				Status:   http.StatusCreated,
				PostForm: `{ "name": "my-stack", "description": "Here is my stack!" }`,
			},
			Expected: stack,
			Action: func(c *Client, t *testing.T) (any, error) {
				stackDescription := "Here is my stack!"
				r := &resource.StackCreate{
					Name:        "my-stack",
					Description: &stackDescription,
				}
				return c.Stacks.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete stack",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/stacks/88db2b75-671f-4e4b-a19a-7db992366595",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Stacks.Delete(context.Background(), "88db2b75-671f-4e4b-a19a-7db992366595")
			},
		},
		{
			Description: "Get stack",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/stacks/88db2b75-671f-4e4b-a19a-7db992366595",
				Output:   g.Single(stack),
				Status:   http.StatusOK},
			Expected: stack,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Stacks.Get(context.Background(), "88db2b75-671f-4e4b-a19a-7db992366595")
			},
		},
		{
			Description: "List all stacks",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/stacks",
				Output:   g.Paged([]string{stack}, []string{stack2}),
				Status:   http.StatusOK},
			Expected: g.Array(stack, stack2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Stacks.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List all apps for given stack",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/stacks/88db2b75-671f-4e4b-a19a-7db992366595/apps",
				Output:   g.Paged([]string{app}, []string{app2}),
				Status:   http.StatusOK},
			Expected: g.Array(app, app2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Stacks.ListAppsOnStackAll(context.Background(), "88db2b75-671f-4e4b-a19a-7db992366595", nil)
			},
		},
		{
			Description: "Update stack",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/stacks/88db2b75-671f-4e4b-a19a-7db992366595",
				Output:   g.Single(stack),
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": { "key": "value" }, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: stack,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.StackUpdate{
					Metadata: resource.NewMetadata().
						WithLabel("", "key", "value").
						WithAnnotation("", "note", "detailed information"),
				}
				return c.Stacks.Update(context.Background(), "88db2b75-671f-4e4b-a19a-7db992366595", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
