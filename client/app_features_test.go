package client

import (
	"github.com/cloudfoundry-community/go-cfclient/test"
	"net/http"
	"testing"
)

func TestAppFeatures(t *testing.T) {
	g := test.NewObjectJSONGenerator(163)
	appFeature := g.AppFeature()

	tests := []RouteTest{
		{
			Description: "Get SSH app feature",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/features/ssh",
				Output:   []string{appFeature},
				Status:   http.StatusOK},
			Expected: appFeature,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AppFeatures.GetSSH("1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "List all app features",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/features",
				Output:   g.Paged([]string{appFeature}),
				Status:   http.StatusOK},
			Expected: g.Array(appFeature),
			Action: func(c *Client, t *testing.T) (any, error) {
				f, _, err := c.AppFeatures.List("1cb006ee-fb05-47e1-b541-c34179ddc446")
				return f, err
			},
		},
		{
			Description: "Update SSH app feature",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/features/ssh",
				Output:   []string{appFeature},
				Status:   http.StatusOK,
				PostForm: `{ "enabled": false }`,
			},
			Expected: appFeature,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AppFeatures.UpdateSSH("1cb006ee-fb05-47e1-b541-c34179ddc446", false)
			},
		},
	}
	executeTests(tests, t)
}
