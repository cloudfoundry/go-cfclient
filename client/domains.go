package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type DomainClient commonClient

// DomainListOptions list filters
type DomainListOptions struct {
	*ListOptions

	GUIDs             Filter `filter:"guids,omitempty"`
	Names             Filter `filter:"names,omitempty"`
	OrganizationGUIDs Filter `filter:"organization_guids,omitempty"`
}

// NewDomainListOptions creates new options to pass to list
func NewDomainListOptions() *DomainListOptions {
	return &DomainListOptions{
		ListOptions: NewListOptions(),
	}
}

// Create a new domain
func (c *DomainClient) Create(r *resource.DomainCreate) (*resource.Domain, error) {
	var d resource.Domain
	err := c.client.post(r.Name, "/v3/domains", r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Delete the specified app
func (c *DomainClient) Delete(guid string) error {
	return c.client.delete(path("/v3/domains/%s", guid))
}

// Get the specified domain
func (c *DomainClient) Get(guid string) (*resource.Domain, error) {
	var d resource.Domain
	err := c.client.get(path("/v3/domains/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// List all Domains the user has access to in paged results
func (c *DomainClient) List(opts *DomainListOptions) ([]*resource.Domain, *Pager, error) {
	var res resource.DomainList
	err := c.client.get(path("/v3/domains?%s", opts.ToQueryString(opts)), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all Domains the user has access to
func (c *DomainClient) ListAll() ([]*resource.Domain, error) {
	opts := NewDomainListOptions()
	var allDomains []*resource.Domain
	for {
		Domains, pager, err := c.List(opts)
		if err != nil {
			return nil, err
		}
		allDomains = append(allDomains, Domains...)
		if !pager.HasNextPage() {
			break
		}
		opts.ListOptions = pager.NextPage(opts.ListOptions)
	}
	return allDomains, nil
}

// ListForOrg retrieves all domains for the specified org that the user has access to in paged results
func (c *DomainClient) ListForOrg(orgGUID string, opts *DomainListOptions) ([]*resource.Domain, *Pager, error) {
	var res resource.DomainList
	err := c.client.get(path("/v3/organizations/%s/domains?%s", orgGUID, opts.ToQueryString(opts)), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// Share an organization-scoped domain to the organization specified by the org guid
// This will allow the organization to use the organization-scoped domain
func (c *DomainClient) Share(domainGUID, orgGUID string) (*resource.ToManyRelationships, error) {
	r := resource.NewDomainShare(orgGUID)
	return c.ShareMany(domainGUID, r)
}

// ShareMany shares an organization-scoped domain to other organizations specified by a list of organization guids
// This will allow any of the other organizations to use the organization-scoped domain.
func (c *DomainClient) ShareMany(guid string, r *resource.ToManyRelationships) (*resource.ToManyRelationships, error) {
	var d resource.ToManyRelationships
	err := c.client.post(guid, path("/v3/domains/%s/relationships/shared_organizations", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Share an organization-scoped domain to other organizations specified by a list of organization guids
// This will allow any of the other organizations to use the organization-scoped domain.
func (c *DomainClient) Unshare(domainGUID, orgGUID string) error {
	return c.client.delete(path("/v3/domains/%s/relationships/shared_organizations/%s", domainGUID, orgGUID))
}

// Update the specified attributes of the app
func (c *DomainClient) Update(guid string, r *resource.DomainUpdate) (*resource.Domain, error) {
	var d resource.Domain
	err := c.client.patch(path("/v3/domains/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
