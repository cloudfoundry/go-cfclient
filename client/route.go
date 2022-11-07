package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type RouteClient commonClient

// RouteListOptions list filters
type RouteListOptions struct {
	*ListOptions

	AppGUIDs             Filter `filter:"app_guids,omitempty"`
	SpaceGUIDs           Filter `filter:"space_guids,omitempty"`
	DomainGUIDs          Filter `filter:"domain_guids,omitempty"`
	OrganizationGUIDs    Filter `filter:"organization_guids,omitempty"`
	ServiceInstanceGUIDs Filter `filter:"service_instance_guids,omitempty"`

	Hosts Filter `filter:"hosts,omitempty"`
	Paths Filter `filter:"paths,omitempty"`
	Ports Filter `filter:"ports,omitempty"`

	Include resource.RouteIncludeType `filter:"include,omitempty"`
}

// NewRouteListOptions creates new options to pass to list
func NewRouteListOptions() *RouteListOptions {
	return &RouteListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o RouteListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new route
func (c *RouteClient) Create(r *resource.RouteCreate) (*resource.Route, error) {
	var Route resource.Route
	_, err := c.client.post("/v3/routes", r, &Route)
	if err != nil {
		return nil, err
	}
	return &Route, nil
}

// Delete the specified route
func (c *RouteClient) Delete(guid string) error {
	return c.client.delete(path("/v3/routes/%s", guid))
}

// Get the specified route
func (c *RouteClient) Get(guid string) (*resource.Route, error) {
	var Route resource.Route
	err := c.client.get(path("/v3/routes/%s", guid), &Route)
	if err != nil {
		return nil, err
	}
	return &Route, nil
}

// GetIncludeDomain allows callers to fetch a route and include the parent domain
func (c *RouteClient) GetIncludeDomain(guid string) (*resource.Route, *resource.Domain, error) {
	var r resource.RouteWithIncluded
	err := c.client.get(path("/v3/routes/%s?include=%s", guid, resource.RouteIncludeDomain), &r)
	if err != nil {
		return nil, nil, err
	}
	return &r.Route, r.Included.Domains[0], nil
}

// GetIncludeSpace allows callers to fetch a route and include the parent space
func (c *RouteClient) GetIncludeSpace(guid string) (*resource.Route, *resource.Space, error) {
	var r resource.RouteWithIncluded
	err := c.client.get(path("/v3/routes/%s?include=%s", guid, resource.RouteIncludeSpaceOrganization), &r)
	if err != nil {
		return nil, nil, err
	}
	return &r.Route, r.Included.Spaces[0], nil
}

// GetIncludeSpaceAndOrg allows callers to fetch a route and include the parent space and org
func (c *RouteClient) GetIncludeSpaceAndOrg(guid string) (*resource.Route, *resource.Space, *resource.Organization, error) {
	var r resource.RouteWithIncluded
	err := c.client.get(path("/v3/routes/%s?include=%s", guid, resource.RouteIncludeSpaceOrganization), &r)
	if err != nil {
		return nil, nil, nil, err
	}
	return &r.Route, r.Included.Spaces[0], r.Included.Organizations[0], nil
}

// List pages routes the user has access to
func (c *RouteClient) List(opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	opts.Include = resource.RouteIncludeNone

	var res resource.RouteList
	err := c.client.get(path("/v3/routes?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all routes the user has access to
func (c *RouteClient) ListAll(opts *RouteListOptions) ([]*resource.Route, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	return AutoPage[*RouteListOptions, *resource.Route](opts, func(opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
		return c.List(opts)
	})
}

// ListForApp pages routes for the specified app the user has access to
func (c *RouteClient) ListForApp(appGUID string, opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	opts.Include = resource.RouteIncludeNone

	var res resource.RouteList
	err := c.client.get(path("/v3/apps/%s/routes?%s", appGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppAll retrieves all routes for the specified app the user has access to
func (c *RouteClient) ListForAppAll(appGUID string, opts *RouteListOptions) ([]*resource.Route, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	return AutoPage[*RouteListOptions, *resource.Route](opts, func(opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
		return c.ListForApp(appGUID, opts)
	})
}

// ListIncludeDomains page all routes the user has access to and include the parent domains
func (c *RouteClient) ListIncludeDomains(opts *RouteListOptions) ([]*resource.Route, []*resource.Domain, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	opts.Include = resource.RouteIncludeDomain

	var res resource.RouteList
	err := c.client.get(path("/v3/routes?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Domains, pager, nil
}

// ListIncludeDomainsAll retrieves all routes the user has access to and includes the parent domains
func (c *RouteClient) ListIncludeDomainsAll(opts *RouteListOptions) ([]*resource.Route, []*resource.Domain, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}

	var all []*resource.Route
	var allDomains []*resource.Domain
	for {
		page, domains, pager, err := c.ListIncludeDomains(opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allDomains = append(allDomains, domains...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allDomains, nil
}

// ListIncludeSpaces page all routes the user has access to and include the parent spaces
func (c *RouteClient) ListIncludeSpaces(opts *RouteListOptions) ([]*resource.Route, []*resource.Space, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	opts.Include = resource.RouteIncludeSpace

	var res resource.RouteList
	err := c.client.get(path("/v3/routes?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Spaces, pager, nil
}

// ListIncludeSpacesAll retrieves all routes the user has access to and includes the parent spaces
func (c *RouteClient) ListIncludeSpacesAll(opts *RouteListOptions) ([]*resource.Route, []*resource.Space, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}

	var all []*resource.Route
	var allSpaces []*resource.Space
	for {
		page, spaces, pager, err := c.ListIncludeSpaces(opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allSpaces = append(allSpaces, spaces...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allSpaces, nil
}

// ListIncludeSpacesAndOrgs page all routes the user has access to and include the parent spaces and orgs
func (c *RouteClient) ListIncludeSpacesAndOrgs(opts *RouteListOptions) ([]*resource.Route, []*resource.Space, []*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	opts.Include = resource.RouteIncludeSpaceOrganization

	var res resource.RouteList
	err := c.client.get(path("/v3/routes?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Spaces, res.Included.Organizations, pager, nil
}

// ListIncludeSpacesAndOrgsAll retrieves all routes the user has access to and includes the parent spaces and org
func (c *RouteClient) ListIncludeSpacesAndOrgsAll(opts *RouteListOptions) ([]*resource.Route, []*resource.Space, []*resource.Organization, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}

	var all []*resource.Route
	var allSpaces []*resource.Space
	var allOrgs []*resource.Organization
	for {
		page, spaces, orgs, pager, err := c.ListIncludeSpacesAndOrgs(opts)
		if err != nil {
			return nil, nil, nil, err
		}
		all = append(all, page...)
		allSpaces = append(allSpaces, spaces...)
		allOrgs = append(allOrgs, orgs...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allSpaces, allOrgs, nil
}

// Update the specified attributes of the app
func (c *RouteClient) Update(guid string, r *resource.RouteUpdate) (*resource.Route, error) {
	var res resource.Route
	err := c.client.patch(path("/v3/routes/%s", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
