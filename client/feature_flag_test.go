package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/cloudfoundry/go-cfclient/v3/testutil"
)

func TestFeatureFlags(t *testing.T) {
	g := testutil.NewObjectJSONGenerator()
	ff := g.FeatureFlag().JSON
	ff2 := g.FeatureFlag().JSON
	ff3 := g.FeatureFlag().JSON
	ff4 := g.FeatureFlag().JSON

	tests := []RouteTest{
		{
			Description: "Get feature flag",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/feature_flags/resource_matching",
				Output:   g.Single(ff),
				Status:   http.StatusOK},
			Expected: ff,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.FeatureFlags.Get(context.Background(), resource.FeatureFlagResourceMatching)
			},
		},
		{
			Description: "List all feature flags",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/feature_flags",
				Output:   g.Paged([]string{ff, ff2}, []string{ff3, ff4}),
				Status:   http.StatusOK},
			Expected: g.Array(ff, ff2, ff3, ff4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.FeatureFlags.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "Update feature flag",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/feature_flags/resource_matching",
				Output:   g.Single(ff),
				Status:   http.StatusOK,
				PostForm: `{ "enabled": true, "custom_error_message": "error message the user sees" }`,
			},
			Expected: ff,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewFeatureFlagUpdate().
					WithEnabled(true).
					WithCustomErrorMessage("error message the user sees")
				return c.FeatureFlags.Update(context.Background(), resource.FeatureFlagResourceMatching, r)
			},
		},
	}
	ExecuteTests(tests, t)
}
