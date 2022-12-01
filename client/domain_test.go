package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestDomains(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	domain := g.Domain().JSON
	domain2 := g.Domain().JSON
	domain3 := g.Domain().JSON
	domain4 := g.Domain().JSON
	sharedDomains := g.DomainShared().JSON

	tests := []RouteTest{
		{
			Description: "Create domain",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/domains",
				Output:   g.Single(domain),
				Status:   http.StatusCreated,
				PostForm: `{ "name": "foo.example.org", "internal": true }`,
			},
			Expected: domain,
			Action: func(c *Client, t *testing.T) (any, error) {
				internal := true
				r := resource.NewDomainCreate("foo.example.org")
				r.Internal = &internal
				return c.Domains.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete domain",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/domains/f666ffc5-106e-4fda-b56f-568b5cf3ae9f",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/api/v3/jobs/c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			},
			Expected: "c33a5caf-77e0-4d6e-b587-5555d339bc9a",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Domains.Delete(context.Background(), "f666ffc5-106e-4fda-b56f-568b5cf3ae9f")
			},
		},
		{
			Description: "Get domain",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/domains/f666ffc5-106e-4fda-b56f-568b5cf3ae9f",
				Output:   g.Single(domain),
				Status:   http.StatusOK},
			Expected: domain,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Domains.Get(context.Background(), "f666ffc5-106e-4fda-b56f-568b5cf3ae9f")
			},
		},
		{
			Description: "List first page of domains",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/domains",
				Output:   g.SinglePaged(domain),
				Status:   http.StatusOK,
			},
			Expected: g.Array(domain),
			Action: func(c *Client, t *testing.T) (any, error) {
				apps, _, err := c.Domains.List(context.Background(), NewDomainListOptions())
				return apps, err
			},
		},
		{
			Description: "List all domains",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/domains",
				Output:   g.Paged([]string{domain, domain2}, []string{domain3, domain4}),
				Status:   http.StatusOK},
			Expected: g.Array(domain, domain2, domain3, domain4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Domains.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "List first page of domains for organization",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3a5f687b-2ce8-4ade-be75-8eca99b0db8b/domains",
				Output:   g.SinglePaged(domain),
				Status:   http.StatusOK,
			},
			Expected: g.Array(domain),
			Action: func(c *Client, t *testing.T) (any, error) {
				apps, _, err := c.Domains.ListForOrganization(context.Background(), "3a5f687b-2ce8-4ade-be75-8eca99b0db8b", NewDomainListOptions())
				return apps, err
			},
		},
		{
			Description: "Update domain",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/domains/f666ffc5-106e-4fda-b56f-568b5cf3ae9f",
				Output:   g.Single(domain),
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": {"key": "value"}, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: domain,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.DomainUpdate{
					Metadata: resource.NewMetadata().
						WithLabel("", "key", "value").
						WithAnnotation("", "note", "detailed information"),
				}
				return c.Domains.Update(context.Background(), "f666ffc5-106e-4fda-b56f-568b5cf3ae9f", r)
			},
		},
		{
			Description: "Share domain",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/domains/1cb006ee-fb05-47e1-b541-c34179ddc446/relationships/shared_organizations",
				Output:   g.Single(sharedDomains),
				Status:   http.StatusOK,
			},
			Expected: sharedDomains,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Domains.Share(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446", "3a5f687b-2ce8-4ade-be75-8eca99b0db8b")
			},
		},
		{
			Description: "Un-share domain",
			Route: testutil.MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/domains/1cb006ee-fb05-47e1-b541-c34179ddc446/relationships/shared_organizations/3a5f687b-2ce8-4ade-be75-8eca99b0db8b",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Domains.UnShare(context.Background(), "1cb006ee-fb05-47e1-b541-c34179ddc446", "3a5f687b-2ce8-4ade-be75-8eca99b0db8b")
			},
		},
	}
	ExecuteTests(tests, t)
}
