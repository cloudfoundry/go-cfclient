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

// RouteListIncludeOptions list filters
type RouteListIncludeOptions struct {
	*RouteListOptions

	Include resource.RouteIncludeType `filter:"include,omitempty"`
}

// NewRouteListIncludeOptions creates new options to pass to list
func NewRouteListIncludeOptions(include resource.RouteIncludeType) *RouteListIncludeOptions {
	return &RouteListIncludeOptions{
		Include:          include,
		RouteListOptions: NewRouteListOptions(),
	}
}

func (o RouteListIncludeOptions) ToQueryString() url.Values {
	u := o.RouteListOptions.ToQueryString()
	if o.Include != resource.RouteIncludeNone {
		u.Set("include", o.Include.String())
	}
	return u
}

// Create a new route
func (c *RouteClient) Create(r *resource.RouteCreate) (*resource.Route, error) {
	var Route resource.Route
	err := c.client.post("", "/v3/routes", r, &Route)
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

// GetInclude allows callers to fetch a route and include information of parent objects in the response
func (c *RouteClient) GetInclude(guid string, include resource.RouteIncludeType) (*resource.Route, *resource.RouteIncluded, error) {
	var r resource.RouteWithIncluded
	err := c.client.get(path("/v3/routes/%s?include=%s", guid, include), &r)
	if err != nil {
		return nil, nil, err
	}
	return &r.Route, r.Included, nil
}

// List pages routes the user has access to
func (c *RouteClient) List(opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
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

// ListInclude page all routes the user has access to and include the specified parent resources
func (c *RouteClient) ListInclude(opts *RouteListIncludeOptions) ([]*resource.Route, *resource.RouteIncluded, *Pager, error) {
	if opts == nil {
		opts = NewRouteListIncludeOptions(resource.RouteIncludeNone)
	}
	var res resource.RouteList
	err := c.client.get(path("/v3/routes?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included, pager, nil
}

// ListIncludeAll retrieves all routes the user has access to and include the specified parent resources
func (c *RouteClient) ListIncludeAll(opts *RouteListIncludeOptions) ([]*resource.Route, *resource.RouteIncluded, error) {
	if opts == nil {
		opts = NewRouteListIncludeOptions(resource.RouteIncludeNone)
	}
	return routeAutoPageInclude[*RouteListIncludeOptions, *resource.Route](opts, func(opts *RouteListIncludeOptions) ([]*resource.Route, *resource.RouteIncluded, *Pager, error) {
		return c.ListInclude(opts)
	})
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

type routeListIncludeFunc[T ListOptioner, R any] func(opts T) ([]R, *resource.RouteIncluded, *Pager, error)

func routeAutoPageInclude[T ListOptioner, R any](opts T, list routeListIncludeFunc[T, R]) ([]R, *resource.RouteIncluded, error) {
	var all []R
	var allIncluded *resource.RouteIncluded
	for {
		page, included, pager, err := list(opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allIncluded.Organizations = append(allIncluded.Organizations, included.Organizations...)
		allIncluded.Spaces = append(allIncluded.Spaces, included.Spaces...)
		allIncluded.Domains = append(allIncluded.Domains, included.Domains...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allIncluded, nil
}
