package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type OrgQuotaClient commonClient

// OrgQuotaListOptions list filters
type OrgQuotaListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`
	Names             Filter `qs:"names"`
	OrganizationGUIDs Filter `qs:"organization_guids"`
}

// NewOrgQuotaListOptions creates new options to pass to list
func NewOrgQuotaListOptions() *OrgQuotaListOptions {
	return &OrgQuotaListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o OrgQuotaListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Apply the specified org quota to the orgs
func (c *OrgQuotaClient) Apply(ctx context.Context, guid string, orgGUIDs []string) ([]string, error) {
	req := resource.NewToManyRelationships(orgGUIDs)
	var relation resource.ToManyRelationships
	_, err := c.client.post(ctx, path.Format("/v3/organization_quotas/%s/relationships/organizations", guid), req, &relation)
	if err != nil {
		return nil, err
	}
	var guids []string
	for _, r := range relation.Data {
		guids = append(guids, r.GUID)
	}
	return guids, nil
}

// Create a new org quota
func (c *OrgQuotaClient) Create(ctx context.Context, r *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error) {
	var q resource.OrganizationQuota
	_, err := c.client.post(ctx, "/v3/organization_quotas", r, &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// Delete the specified org quota
func (c *OrgQuotaClient) Delete(ctx context.Context, guid string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/organization_quotas/%s", guid))
	return err
}

// First returns the first organization quota matching the options or an error when less than 1 match
func (c *OrgQuotaClient) First(ctx context.Context, opts *OrgQuotaListOptions) (*resource.OrganizationQuota, error) {
	return First[*OrgQuotaListOptions, *resource.OrganizationQuota](opts, func(opts *OrgQuotaListOptions) ([]*resource.OrganizationQuota, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified org quota
func (c *OrgQuotaClient) Get(ctx context.Context, guid string) (*resource.OrganizationQuota, error) {
	var app resource.OrganizationQuota
	err := c.client.get(ctx, path.Format("/v3/organization_quotas/%s", guid), &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// List pages all org quotas the user has access to
func (c *OrgQuotaClient) List(ctx context.Context, opts *OrgQuotaListOptions) ([]*resource.OrganizationQuota, *Pager, error) {
	if opts == nil {
		opts = NewOrgQuotaListOptions()
	}

	var res resource.OrganizationQuotaList
	err := c.client.get(ctx, path.Format("/v3/organization_quotas?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all org quotas the user has access to
func (c *OrgQuotaClient) ListAll(ctx context.Context, opts *OrgQuotaListOptions) ([]*resource.OrganizationQuota, error) {
	if opts == nil {
		opts = NewOrgQuotaListOptions()
	}
	return AutoPage[*OrgQuotaListOptions, *resource.OrganizationQuota](opts, func(opts *OrgQuotaListOptions) ([]*resource.OrganizationQuota, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Single returns a single org quota matching the options or an error if not exactly 1 match
func (c *OrgQuotaClient) Single(ctx context.Context, opts *OrgQuotaListOptions) (*resource.OrganizationQuota, error) {
	return Single[*OrgQuotaListOptions, *resource.OrganizationQuota](opts, func(opts *OrgQuotaListOptions) ([]*resource.OrganizationQuota, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the org quota
func (c *OrgQuotaClient) Update(ctx context.Context, guid string, r *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error) {
	var q resource.OrganizationQuota
	_, err := c.client.patch(ctx, path.Format("/v3/organization_quotas/%s", guid), r, &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
