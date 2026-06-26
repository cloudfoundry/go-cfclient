package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

type RoutePolicyClient commonClient

// RoutePolicyListOptions list filters
type RoutePolicyListOptions struct {
	*ListOptions

	GUIDs             Filter                          `qs:"guids"`
	RouteGUIDs        Filter                          `qs:"route_guids"`
	SpaceGUIDs        Filter                          `qs:"space_guids"`
	SourceGUIDs       Filter                          `qs:"source_guids"`
	Sources           Filter                          `qs:"sources"`
	OrganizationGUIDs Filter                          `qs:"organization_guids"`
	Include           resource.RoutePolicyIncludeType `qs:"include"`
}

// NewRoutePolicyListOptions creates new options to pass to list
func NewRoutePolicyListOptions() *RoutePolicyListOptions {
	return &RoutePolicyListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o RoutePolicyListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new route policy
func (c *RoutePolicyClient) Create(ctx context.Context, r *resource.RoutePolicyCreate) (*resource.RoutePolicy, error) {
	var routePolicy resource.RoutePolicy
	_, err := c.client.post(ctx, "/v3/route_policies", r, &routePolicy)
	if err != nil {
		return nil, err
	}
	return &routePolicy, nil
}

// Delete the specified route policy asynchronously and return a jobGUID.
func (c *RoutePolicyClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/route_policies/%s", guid))
}

// First returns the first route policy matching the options or an error when less than 1 match
func (c *RoutePolicyClient) First(ctx context.Context, opts *RoutePolicyListOptions) (*resource.RoutePolicy, error) {
	return First[*RoutePolicyListOptions, *resource.RoutePolicy](opts, func(opts *RoutePolicyListOptions) ([]*resource.RoutePolicy, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified route policy
func (c *RoutePolicyClient) Get(ctx context.Context, guid string) (*resource.RoutePolicy, error) {
	var routePolicy resource.RoutePolicy
	err := c.client.get(ctx, path.Format("/v3/route_policies/%s", guid), &routePolicy)
	if err != nil {
		return nil, err
	}
	return &routePolicy, nil
}

// GetIncludeRoute retrieves the specified route policy and include the associated route
func (c *RoutePolicyClient) GetIncludeRoute(ctx context.Context, guid string) (*resource.RoutePolicy, *resource.Route, error) {
	var res resource.RoutePolicyList
	err := c.client.get(ctx, path.Format("/v3/route_policies/%s?include=%s", guid, resource.RoutePolicyIncludeRoute), &res)
	if err != nil {
		return nil, nil, err
	}
	if len(res.Resources) == 0 {
		return nil, nil, err
	}
	var route *resource.Route
	if res.Included != nil && len(res.Included.Routes) > 0 {
		route = res.Included.Routes[0]
	}
	return res.Resources[0], route, nil
}

// GetIncludeSource retrieves the specified route policy and include the associated source app
func (c *RoutePolicyClient) GetIncludeSource(ctx context.Context, guid string) (*resource.RoutePolicy, *resource.App, error) {
	var res resource.RoutePolicyList
	err := c.client.get(ctx, path.Format("/v3/route_policies/%s?include=%s", guid, resource.RoutePolicyIncludeSource), &res)
	if err != nil {
		return nil, nil, err
	}
	if len(res.Resources) == 0 {
		return nil, nil, err
	}
	var app *resource.App
	if res.Included != nil && len(res.Included.Apps) > 0 {
		app = res.Included.Apps[0]
	}
	return res.Resources[0], app, nil
}

// List pages all the route policies the user has access to
func (c *RoutePolicyClient) List(ctx context.Context, opts *RoutePolicyListOptions) ([]*resource.RoutePolicy, *Pager, error) {
	if opts == nil {
		opts = NewRoutePolicyListOptions()
	}
	opts.Include = resource.RoutePolicyIncludeNone

	var res resource.RoutePolicyList
	err := c.client.list(ctx, "/v3/route_policies", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all route policies the user has access to
func (c *RoutePolicyClient) ListAll(ctx context.Context, opts *RoutePolicyListOptions) ([]*resource.RoutePolicy, error) {
	if opts == nil {
		opts = NewRoutePolicyListOptions()
	}
	return AutoPage[*RoutePolicyListOptions, *resource.RoutePolicy](opts, func(opts *RoutePolicyListOptions) ([]*resource.RoutePolicy, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListIncludeRoute page all route policies the user has access to and include the associated routes
func (c *RoutePolicyClient) ListIncludeRoute(ctx context.Context, opts *RoutePolicyListOptions) ([]*resource.RoutePolicy, []*resource.Route, *Pager, error) {
	if opts == nil {
		opts = NewRoutePolicyListOptions()
	}
	opts.Include = resource.RoutePolicyIncludeRoute

	var res resource.RoutePolicyList
	err := c.client.list(ctx, "/v3/route_policies", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Routes, pager, nil
}

// ListIncludeRouteAll retrieves all route policies the user has access to and include the associated routes
func (c *RoutePolicyClient) ListIncludeRouteAll(ctx context.Context, opts *RoutePolicyListOptions) ([]*resource.RoutePolicy, []*resource.Route, error) {
	if opts == nil {
		opts = NewRoutePolicyListOptions()
	}

	var all []*resource.RoutePolicy
	var allRoutes []*resource.Route
	for {
		page, routes, pager, err := c.ListIncludeRoute(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allRoutes = append(allRoutes, routes...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allRoutes, nil
}

// ListIncludeSource page all route policies the user has access to and include the associated source apps
func (c *RoutePolicyClient) ListIncludeSource(ctx context.Context, opts *RoutePolicyListOptions) ([]*resource.RoutePolicy, []*resource.App, []*resource.Space, []*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewRoutePolicyListOptions()
	}
	opts.Include = resource.RoutePolicyIncludeSource

	var res resource.RoutePolicyList
	err := c.client.list(ctx, "/v3/route_policies", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Apps, res.Included.Spaces, res.Included.Organizations, pager, nil
}

// ListIncludeSourceAll retrieves all route policies the user has access to and include the associated source apps
func (c *RoutePolicyClient) ListIncludeSourceAll(ctx context.Context, opts *RoutePolicyListOptions) ([]*resource.RoutePolicy, []*resource.App, []*resource.Space, []*resource.Organization, error) {
	if opts == nil {
		opts = NewRoutePolicyListOptions()
	}

	var all []*resource.RoutePolicy
	var allApps []*resource.App
	var allSpaces []*resource.Space
	var allOrganizations []*resource.Organization
	for {
		page, apps, spaces, orgs, pager, err := c.ListIncludeSource(ctx, opts)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		all = append(all, page...)
		allApps = append(allApps, apps...)
		allSpaces = append(allSpaces, spaces...)
		allOrganizations = append(allOrganizations, orgs...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allApps, allSpaces, allOrganizations, nil
}

// Single returns a single route policy matching the options or an error if not exactly 1 match
func (c *RoutePolicyClient) Single(ctx context.Context, opts *RoutePolicyListOptions) (*resource.RoutePolicy, error) {
	return Single[*RoutePolicyListOptions, *resource.RoutePolicy](opts, func(opts *RoutePolicyListOptions) ([]*resource.RoutePolicy, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the route policy
func (c *RoutePolicyClient) Update(ctx context.Context, guid string, r *resource.RoutePolicyUpdate) (*resource.RoutePolicy, error) {
	var routePolicy resource.RoutePolicy
	_, err := c.client.patch(ctx, path.Format("/v3/route_policies/%s", guid), r, &routePolicy)
	if err != nil {
		return nil, err
	}
	return &routePolicy, nil
}
