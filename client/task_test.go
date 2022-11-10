package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestTasks(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	task := g.Task()
	task2 := g.Task()
	task3 := g.Task()
	task4 := g.Task()

	tests := []RouteTest{
		{
			Description: "Cancel a task",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/tasks/d9442132-4669-49f7-a3c5-8fa8d1150504/actions/cancel",
				Output:   []string{task},
				Status:   http.StatusOK,
			},
			Expected: task,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Tasks.Cancel("d9442132-4669-49f7-a3c5-8fa8d1150504")
			},
		},
		{
			Description: "Create a task",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/apps/631b46a1-c3b6-4599-9659-72c9fd54817f/tasks",
				Output:   []string{task},
				Status:   http.StatusCreated,
				PostForm: `{ "command": "rake db:migrate" }`,
			},
			Expected: task,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewTaskCreateWithCommand("rake db:migrate")
				return c.Tasks.Create("631b46a1-c3b6-4599-9659-72c9fd54817f", r)
			},
		},
		{
			Description: "Get task",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/tasks/d9442132-4669-49f7-a3c5-8fa8d1150504",
				Output:   []string{task},
				Status:   http.StatusOK},
			Expected: task,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Tasks.Get("d9442132-4669-49f7-a3c5-8fa8d1150504")
			},
		},
		{
			Description: "List all tasks",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/tasks",
				Output:   g.Paged([]string{task, task2}, []string{task3, task4}),
				Status:   http.StatusOK},
			Expected: g.Array(task, task2, task3, task4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Tasks.ListAll(nil)
			},
		},
		{
			Description: "List all tasks for an app",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/631b46a1-c3b6-4599-9659-72c9fd54817f/tasks",
				Output:   g.Paged([]string{task, task2}, []string{task3, task4}),
				Status:   http.StatusOK},
			Expected: g.Array(task, task2, task3, task4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Tasks.ListForAppAll("631b46a1-c3b6-4599-9659-72c9fd54817f", nil)
			},
		},
		{
			Description: "Update a task",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/tasks/d9442132-4669-49f7-a3c5-8fa8d1150504",
				Output:   []string{task},
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": { "key": "value" }, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: task,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.TaskUpdate{
					Metadata: &resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.Tasks.Update("d9442132-4669-49f7-a3c5-8fa8d1150504", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
