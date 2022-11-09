package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"net/http"
	"testing"
)

func TestAuditEvents(t *testing.T) {
	g := test.NewObjectJSONGenerator(161)
	auditEvent := g.AuditEvent()
	auditEvent2 := g.AuditEvent()
	auditEvent3 := g.AuditEvent()

	tests := []RouteTest{
		{
			Description: "Get audit event",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/audit_events/27a9b4a5-ba8a-448c-ac51-3a6dab9aa3f8",
				Output:   []string{auditEvent},
				Status:   http.StatusOK},
			Expected: auditEvent,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AuditEvents.Get("27a9b4a5-ba8a-448c-ac51-3a6dab9aa3f8")
			},
		},
		{
			Description: "List all audit events",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/audit_events",
				Output:   g.Paged([]string{auditEvent, auditEvent2}, []string{auditEvent3}),
				Status:   http.StatusOK},
			Expected: g.Array(auditEvent, auditEvent2, auditEvent3),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.AuditEvents.ListAll(nil)
			},
		},
	}
	executeTests(tests, t)
}
