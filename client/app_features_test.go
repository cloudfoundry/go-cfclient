package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestAppFeatures(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(163)
	appFeature := g.AppFeature().JSON

	tests := []RouteTest{
		{
			Description: "Get SSH app feature",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/features/ssh",
				Output:   g.Single(appFeature),
				Status:   http.StatusOK},
			Expected: appFeature,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AppFeatures.GetSSH(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
			},
		},
		{
			Description: "List all app features",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/features",
				Output:   g.SinglePaged(appFeature),
				Status:   http.StatusOK},
			Expected: g.Array(appFeature),
			Action: func(c *Client, t *testing.T) (any, error) {
				f, _, err := c.AppFeatures.List(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446")
				return f, err
			},
		},
		{
			Description: "Update SSH app feature",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/features/ssh",
				Output:   g.Single(appFeature),
				Status:   http.StatusOK,
				PostForm: `{ "enabled": false }`,
			},
			Expected: appFeature,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AppFeatures.UpdateSSH(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446", false)
			},
		},
	}
	ExecuteTests(tests, t)
}
