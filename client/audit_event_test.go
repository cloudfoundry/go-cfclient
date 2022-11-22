package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"net/http"
	"testing"
)

func TestAuditEvents(t *testing.T) {
	g := testutil.NewObjectJSONGenerator(161)
	auditEvent := g.AuditEvent().JSON
	auditEvent2 := g.AuditEvent().JSON
	auditEvent3 := g.AuditEvent().JSON

	tests := []RouteTest{
		{
			Description: "Get audit event",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/audit_events/27a9b4a5-ba8a-448c-ac51-3a6dab9aa3f8",
				Output:   g.Single(auditEvent),
				Status:   http.StatusOK},
			Expected: auditEvent,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AuditEvents.Get(context.Background(), "27a9b4a5-ba8a-448c-ac51-3a6dab9aa3f8")
			},
		},
		{
			Description: "List all audit events",
			Route: testutil.MockRoute{
				Method:   "GET",
				Endpoint: "/v3/audit_events",
				Output:   g.Paged([]string{auditEvent, auditEvent2}, []string{auditEvent3}),
				Status:   http.StatusOK},
			Expected: g.Array(auditEvent, auditEvent2, auditEvent3),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AuditEvents.ListAll(context.Background(), nil)
			},
		},
	}
	ExecuteTests(tests, t)
}
