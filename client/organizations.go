package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"net/url"
)

type OrgClient commonClient

type OrgListOptions struct {
	*ListOptions

	GUIDs Filter `filter:"guids,omitempty"` // list of organization guids to filter by
	Names Filter `filter:"names,omitempty"` // list of organization names to filter by
}

// NewOrgListOptions creates new options to pass to list
func NewOrgListOptions() *OrgListOptions {
	return &OrgListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o OrgListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create an organization
func (c *OrgClient) Create(r *resource.OrganizationCreate) (*resource.Organization, error) {
	var org resource.Organization
	err := c.client.post(r.Name, "/v3/organizations", r, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// Delete the specified organization
func (c *OrgClient) Delete(guid string) error {
	return c.client.delete(path("/v3/organizations/%s", guid))
}

// Get the specified organization
func (c *OrgClient) Get(guid string) (*resource.Organization, error) {
	var org resource.Organization
	err := c.client.get(path("/v3/organizations/%s", guid), &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// List pages all organizations the user has access to
func (c *OrgClient) List(opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewOrgListOptions()
	}
	var res resource.OrganizationList
	err := c.client.get(path("/v3/organizations?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all organizations the user has access to
func (c *OrgClient) ListAll(opts *OrgListOptions) ([]*resource.Organization, error) {
	if opts == nil {
		opts = NewOrgListOptions()
	}
	return AutoPage[*OrgListOptions, *resource.Organization](opts, func(opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
		return c.List(opts)
	})
}

// Update the organization's specified attributes
func (c *OrgClient) Update(guid string, r *resource.OrganizationUpdate) (*resource.Organization, error) {
	var org resource.Organization
	err := c.client.patch(path("/v3/organizations/%s", guid), r, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}
