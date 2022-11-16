package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type AuditEventClient commonClient

// AuditEventListOptions list filters
type AuditEventListOptions struct {
	*ListOptions

	Types             Filter `qs:"types"`        //  list of event types to filter by
	TargetGUIDs       Filter `qs:"target_guids"` // list of target guids to filter by
	OrganizationGUIDs Filter `qs:"organization_guids"`
	SpaceGUIDs        Filter `qs:"space_guids"`
}

// NewAuditEventListOptions creates new options to pass to list
func NewAuditEventListOptions() *AuditEventListOptions {
	return &AuditEventListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o AuditEventListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Get retrieves the specified audit event
func (c *AuditEventClient) Get(guid string) (*resource.AuditEvent, error) {
	var a resource.AuditEvent
	err := c.client.get(path.Format("/v3/audit_events/%s", guid), &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// List pages all audit events the user has access to
func (c *AuditEventClient) List(opts *AuditEventListOptions) ([]*resource.AuditEvent, *Pager, error) {
	if opts == nil {
		opts = NewAuditEventListOptions()
	}
	var res resource.AuditEventList
	err := c.client.get(path.Format("/v3/audit_events?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all audit events the user has access to
func (c *AuditEventClient) ListAll(opts *AuditEventListOptions) ([]*resource.AuditEvent, error) {
	if opts == nil {
		opts = NewAuditEventListOptions()
	}

	var all []*resource.AuditEvent
	for {
		page, pager, err := c.List(opts)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, nil
}
