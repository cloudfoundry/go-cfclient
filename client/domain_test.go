package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestDomains(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(1)
	domain := g.Domain()
	domain2 := g.Domain()
	domain3 := g.Domain()
	domain4 := g.Domain()
	sharedDomains := g.DomainShared()

	tests := []RouteTest{
		{
			Description: "Create domain",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/domains",
				Output:   []string{domain},
				Status:   http.StatusCreated,
				PostForm: `{ "name": "foo.example.org", "internal": true }`,
			},
			Expected: domain,
			Action: func(c *Client, t *testing.T) (any, error) {
				internal := true
				r := resource.NewDomainCreate("foo.example.org")
				r.Internal = &internal
				return c.Domains.Create(r)
			},
		},
		{
			Description: "Get domain",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/domains/f666ffc5-106e-4fda-b56f-568b5cf3ae9f",
				Output:   []string{domain},
				Status:   http.StatusOK},
			Expected: domain,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Domains.Get("f666ffc5-106e-4fda-b56f-568b5cf3ae9f")
			},
		},
		{
			Description: "List first page of domains",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/domains",
				Output:   g.Paged([]string{domain}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(domain),
			Action: func(c *Client, t *testing.T) (any, error) {
				apps, _, err := c.Domains.List(NewDomainListOptions())
				return apps, err
			},
		},
		{
			Description: "List all domains",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/domains",
				Output:   g.Paged([]string{domain, domain2}, []string{domain3, domain4}),
				Status:   http.StatusOK},
			Expected: g.Array(domain, domain2, domain3, domain4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Domains.ListAll(nil)
			},
		},
		{
			Description: "List first page of domains for org",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3a5f687b-2ce8-4ade-be75-8eca99b0db8b/domains",
				Output:   g.Paged([]string{domain}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(domain),
			Action: func(c *Client, t *testing.T) (any, error) {
				apps, _, err := c.Domains.ListForOrg("3a5f687b-2ce8-4ade-be75-8eca99b0db8b", NewDomainListOptions())
				return apps, err
			},
		},
		{
			Description: "Update domain",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/domains/f666ffc5-106e-4fda-b56f-568b5cf3ae9f",
				Output:   []string{domain},
				Status:   http.StatusOK,
				PostForm: `{ "metadata": { "labels": {"key": "value"}, "annotations": {"note": "detailed information"}}}`,
			},
			Expected: domain,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.DomainUpdate{
					Metadata: &resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.Domains.Update("f666ffc5-106e-4fda-b56f-568b5cf3ae9f", r)
			},
		},
		{
			Description: "Delete domain",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/domains/f666ffc5-106e-4fda-b56f-568b5cf3ae9f",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Domains.Delete("f666ffc5-106e-4fda-b56f-568b5cf3ae9f")
			},
		},
		{
			Description: "Share domain",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/domains/1cb006ee-fb05-47e1-b541-c34179ddc446/relationships/shared_organizations",
				Output:   []string{sharedDomains},
				Status:   http.StatusOK,
			},
			Expected: sharedDomains,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Domains.Share("1cb006ee-fb05-47e1-b541-c34179ddc446", "3a5f687b-2ce8-4ade-be75-8eca99b0db8b")
			},
		},
		{
			Description: "Un-share domain",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/domains/1cb006ee-fb05-47e1-b541-c34179ddc446/relationships/shared_organizations/3a5f687b-2ce8-4ade-be75-8eca99b0db8b",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Domains.Unshare("1cb006ee-fb05-47e1-b541-c34179ddc446", "3a5f687b-2ce8-4ade-be75-8eca99b0db8b")
			},
		},
	}
	ExecuteTests(tests, t)
}
