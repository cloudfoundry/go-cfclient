package client

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestStacks(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	stack := g.Stack()
	stack2 := g.Stack()

	tests := []RouteTest{
		{
			Description: "Create stack",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/stacks",
				Output:   []string{stack},
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
				return c.Stacks.Create(r)
			},
		},
		{
			Description: "Delete stack",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/stacks/88db2b75-671f-4e4b-a19a-7db992366595",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Stacks.Delete("88db2b75-671f-4e4b-a19a-7db992366595")
			},
		},
		{
			Description: "Get stack",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/stacks/88db2b75-671f-4e4b-a19a-7db992366595",
				Output:   []string{stack},
				Status:   http.StatusOK},
			Expected: stack,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Stacks.Get("88db2b75-671f-4e4b-a19a-7db992366595")
			},
		},
		{
			Description: "List all stacks",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/stacks",
				Output:   g.Paged([]string{stack}, []string{stack2}),
				Status:   http.StatusOK},
			Expected: g.Array(stack, stack2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Stacks.ListAll(nil)
			},
		},
		{
			Description: "Update stack",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/stacks/88db2b75-671f-4e4b-a19a-7db992366595",
				Output:   []string{stack},
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": { "key": "value" }, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: stack,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.StackUpdate{
					Metadata: &resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.Stacks.Update("88db2b75-671f-4e4b-a19a-7db992366595", r)
			},
		},
	}
	for _, tt := range tests {
		func() {
			setup(tt.Route, t)
			defer teardown()
			details := fmt.Sprintf("%s %s", tt.Route.Method, tt.Route.Endpoint)
			if tt.Description != "" {
				details = tt.Description + ": " + details
			}

			c, _ := NewTokenConfig(server.URL, "foobar")
			cl, err := New(c)
			require.NoError(t, err, details)

			obj, err := tt.Action(cl, t)
			require.NoError(t, err, details)
			if tt.Expected != "" {
				actual, err := json.Marshal(obj)
				require.NoError(t, err, details)
				require.JSONEq(t, tt.Expected, string(actual), details)
			}
		}()
	}
}
