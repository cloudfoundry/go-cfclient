package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type OrgQuotaClient commonClient

// OrgQuotaListOptions list filters
type OrgQuotaListOptions struct {
	*ListOptions

	GUIDs             Filter `filter:"guids,omitempty"`
	Names             Filter `filter:"names,omitempty"`
	OrganizationGUIDs Filter `filter:"organization_guids,omitempty"`
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
func (c *OrgQuotaClient) Apply(guid string, orgGUIDs []string) ([]string, error) {
	req := resource.NewToManyRelationships(orgGUIDs)
	var relation resource.ToManyRelationships
	_, err := c.client.post(path("/v3/organization_quotas/%s/relationships/organizations", guid), req, &relation)
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
func (c *OrgQuotaClient) Create(r *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error) {
	var q resource.OrganizationQuota
	_, err := c.client.post("/v3/organization_quotas", r, &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// Delete the specified org quota
func (c *OrgQuotaClient) Delete(guid string) error {
	return c.client.delete(path("/v3/organization_quotas/%s", guid))
}

// Get the specified org quota
func (c *OrgQuotaClient) Get(guid string) (*resource.OrganizationQuota, error) {
	var app resource.OrganizationQuota
	err := c.client.get(path("/v3/organization_quotas/%s", guid), &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// List pages all org quotas the user has access to
func (c *OrgQuotaClient) List(opts *OrgQuotaListOptions) ([]*resource.OrganizationQuota, *Pager, error) {
	if opts == nil {
		opts = NewOrgQuotaListOptions()
	}

	var res resource.OrganizationQuotaList
	err := c.client.get(path("/v3/organization_quotas?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all org quotas the user has access to
func (c *OrgQuotaClient) ListAll(opts *OrgQuotaListOptions) ([]*resource.OrganizationQuota, error) {
	if opts == nil {
		opts = NewOrgQuotaListOptions()
	}
	return AutoPage[*OrgQuotaListOptions, *resource.OrganizationQuota](opts, func(opts *OrgQuotaListOptions) ([]*resource.OrganizationQuota, *Pager, error) {
		return c.List(opts)
	})
}

// Update the specified attributes of the org quota
func (c *OrgQuotaClient) Update(guid string, r *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error) {
	var q resource.OrganizationQuota
	err := c.client.patch(path("/v3/organization_quotas/%s", guid), r, &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
