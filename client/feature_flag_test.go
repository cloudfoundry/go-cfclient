package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestFeatureFlags(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(852)
	ff := g.FeatureFlag()
	ff2 := g.FeatureFlag()
	ff3 := g.FeatureFlag()
	ff4 := g.FeatureFlag()

	tests := []RouteTest{
		{
			Description: "Get feature flag",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/feature_flags/resource_matching",
				Output:   []string{ff},
				Status:   http.StatusOK},
			Expected: ff,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.FeatureFlags.Get(resource.FeatureFlagResourceMatching)
			},
		},
		{
			Description: "List all feature flags",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/feature_flags",
				Output:   g.Paged([]string{ff, ff2}, []string{ff3, ff4}),
				Status:   http.StatusOK},
			Expected: g.Array(ff, ff2, ff3, ff4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.FeatureFlags.ListAll(nil)
			},
		},
		{
			Description: "Update feature flag",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/feature_flags/resource_matching",
				Output:   []string{ff},
				Status:   http.StatusOK,
				PostForm: `{ "enabled": true, "custom_error_message": "error message the user sees" }`,
			},
			Expected: ff,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewFeatureFlagUpdate().
					WithEnabled(true).
					WithCustomErrorMessage("error message the user sees")
				return c.FeatureFlags.Update(resource.FeatureFlagResourceMatching, r)
			},
		},
	}
	ExecuteTests(tests, t)
}
