package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestServiceBrokers(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(6872)
	sb := g.ServiceBroker().JSON
	sb2 := g.ServiceBroker().JSON

	tests := []RouteTest{
		{
			Description: "Create service broker",
			Route: testutil.MockRoute{
				Method:   "POST",
				Endpoint: "/v3/service_brokers",
				Output:   g.Single(sb),
				Status:   http.StatusAccepted,
				PostForm: `{
					"name": "my_service_broker",
					"url": "https://example.service-broker.com",
					"authentication": {
					  "type": "basic",
					  "credentials": {
						"username": "us3rn4me",
						"password": "p4ssw0rd"
					  }
					},
					"relationships": {
					  "space": {
						"data": {
						  "guid": "2f35885d-0c9d-4423-83ad-fd05066f8576"
						}
					  }
					}
				  }`,
				RedirectLocation: "https://api.example.org/v3/jobs/af5c57f6-8769-41fa-a499-2c84ed896788",
			},
			Expected: "af5c57f6-8769-41fa-a499-2c84ed896788",
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewServiceBrokerCreate(
					"my_service_broker",
					"https://example.service-broker.com",
					"us3rn4me", "p4ssw0rd")
				r.WithSpace("2f35885d-0c9d-4423-83ad-fd05066f8576")
				return c.ServiceBrokers.Create(context.Background(), r)
			},
		},
		{
			Description: "Delete service broker",
			Route: testutil.MockRoute{
				Method:           "DELETE",
				Endpoint:         "/v3/service_brokers/c680ad12-1ada-4051-8f85-e859e3819c6a",
				Status:           http.StatusAccepted,
				RedirectLocation: "https://api.example.org/v3/jobs/af5c57f6-8769-41fa-a499-2c84ed896788",
			},
			Expected: "af5c57f6-8769-41fa-a499-2c84ed896788",
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceBrokers.Delete(context.Background(), "c680ad12-1ada-4051-8f85-e859e3819c6a")
			},
		},
		{
			Description: "Get service brokers",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_brokers/c680ad12-1ada-4051-8f85-e859e3819c6a",
				Output:   g.Single(sb),
				Status:   http.StatusOK},
			Expected: sb,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceBrokers.Get(context.Background(), "c680ad12-1ada-4051-8f85-e859e3819c6a")
			},
		},
		{
			Description: "List all service brokers",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_brokers",
				Output:   g.Paged([]string{sb}, []string{sb2}),
				Status:   http.StatusOK},
			Expected: g.Array(sb, sb2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServiceBrokers.ListAll(context.Background(), nil)
			},
		},
		{
			Description: "Update a service broker",
			Route: testutil.MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_brokers/c680ad12-1ada-4051-8f85-e859e3819c6a",
				Output:   g.Single(sb),
				Status:   http.StatusAccepted,
				PostForm: `{
					"name": "my_service_broker",
					"url": "https://example.service-broker.com",
					"authentication": {
					  "type": "basic",
					  "credentials": {
						"username": "us3rn4me",
						"password": "p4ssw0rd"
					  }
					}
				  }`,
				RedirectLocation: "https://api.example.org/v3/jobs/af5c57f6-8769-41fa-a499-2c84ed896788",
			},
			Expected: "af5c57f6-8769-41fa-a499-2c84ed896788",
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				r := resource.NewServiceBrokerUpdate().
					WithName("my_service_broker").
					WithURL("https://example.service-broker.com").
					WithCredentials("us3rn4me", "p4ssw0rd")
				return c.ServiceBrokers.Update(context.Background(), "c680ad12-1ada-4051-8f85-e859e3819c6a", r)
			},
		},
	}
	ExecuteTests(tests, t)
}
