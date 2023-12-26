package resource

import (
	"encoding/json"
)

type AuditEvent struct {
	Type string `json:"type"`

	Actor  AuditEventRelatedObject `json:"actor"`
	Target AuditEventRelatedObject `json:"target"`

	Data         *json.RawMessage `json:"data"`
	Space        Relationship     `json:"space"`
	Organization Relationship     `json:"organization"`

	Resource `json:",inline"`
}

type AuditEventList struct {
	Pagination Pagination    `json:"pagination"`
	Resources  []*AuditEvent `json:"resources"`
}

type AuditEventRelatedObject struct {
	GUID string `json:"guid"`
	Type string `json:"type"`
	Name string `json:"name"`
}
