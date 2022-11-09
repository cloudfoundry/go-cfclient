package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type ServiceRouteBindingClient commonClient

// ServiceRouteBindingListOptions list filters
type ServiceRouteBindingListOptions struct {
	*ListOptions

	GUIDs                Filter `filter:"guids,omitempty"`
	RouteGUIDs           Filter `filter:"route_guids,omitempty"`
	ServiceInstanceGUIDs Filter `filter:"service_instance_guids,omitempty"`
	ServiceInstanceNames Filter `filter:"service_instance_names,omitempty"`

	Include resource.ServiceRouteBindingIncludeType `filter:"include,omitempty"`
}

// NewServiceRouteBindingListOptions creates new options to pass to list
func NewServiceRouteBindingListOptions() *ServiceRouteBindingListOptions {
	return &ServiceRouteBindingListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceRouteBindingListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new service route binding
func (c *ServiceRouteBindingClient) Create(r *resource.ServiceRouteBindingCreate) (*resource.ServiceRouteBinding, error) {
	var srb resource.ServiceRouteBinding
	_, err := c.client.post("/v3/service_route_bindings", r, &srb)
	if err != nil {
		return nil, err
	}
	return &srb, nil
}

// Delete the specified service route binding
func (c *ServiceRouteBindingClient) Delete(guid string) error {
	_, err := c.client.delete(path("/v3/service_route_bindings/%s", guid))
	return err
}

// Get the specified service route binding
func (c *ServiceRouteBindingClient) Get(guid string) (*resource.ServiceRouteBinding, error) {
	var srb resource.ServiceRouteBinding
	err := c.client.get(path("/v3/service_route_bindings/%s", guid), &srb)
	if err != nil {
		return nil, err
	}
	return &srb, nil
}

// GetIncludeRoute allows callers to fetch a service route binding and include the associated route
func (c *ServiceRouteBindingClient) GetIncludeRoute(guid string) (*resource.ServiceRouteBinding, *resource.Route, error) {
	var srb resource.ServiceRouteBindingWithIncluded
	err := c.client.get(path("/v3/service_route_bindings/%s?include=%s", guid, resource.ServiceRouteBindingIncludeRoute), &srb)
	if err != nil {
		return nil, nil, err
	}
	return &srb.ServiceRouteBinding, srb.Included.Routes[0], nil
}

// GetIncludeServiceInstance allows callers to fetch a service route binding and include the associated service instance
func (c *ServiceRouteBindingClient) GetIncludeServiceInstance(guid string) (*resource.ServiceRouteBinding, *resource.ServiceInstance, error) {
	var srb resource.ServiceRouteBindingWithIncluded
	err := c.client.get(path("/v3/service_route_bindings/%s?include=%s", guid, resource.ServiceRouteBindingIncludeServiceInstance), &srb)
	if err != nil {
		return nil, nil, err
	}
	return &srb.ServiceRouteBinding, srb.Included.ServiceInstances[0], nil
}

// GetParameters queries the Service Broker for the parameters associated with this service route binding
func (c *ServiceRouteBindingClient) GetParameters(guid string) (map[string]string, error) {
	var srbEnv map[string]string
	err := c.client.get(path("/v3/service_route_bindings/%s/parameters", guid), &srbEnv)
	if err != nil {
		return nil, err
	}
	return srbEnv, nil
}

// List pages all the service route bindings the user has access to
func (c *ServiceRouteBindingClient) List(opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, *Pager, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}
	opts.Include = resource.ServiceRouteBindingIncludeNone

	var res resource.ServiceRouteBindingList
	err := c.client.get(path("/v3/service_route_bindings?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service route bindings the user has access to
func (c *ServiceRouteBindingClient) ListAll(opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}
	return AutoPage[*ServiceRouteBindingListOptions, *resource.ServiceRouteBinding](opts, func(opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, *Pager, error) {
		return c.List(opts)
	})
}

// ListIncludeRoutes page all service route bindings the user has access to and include the associated routes
func (c *ServiceRouteBindingClient) ListIncludeRoutes(opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, []*resource.Route, *Pager, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}
	opts.Include = resource.ServiceRouteBindingIncludeNone

	var res resource.ServiceRouteBindingList
	err := c.client.get(path("/v3/service_route_bindings?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Routes, pager, nil
}

// ListIncludeRoutesAll retrieves all service route bindings the user has access to and include the associated routes
func (c *ServiceRouteBindingClient) ListIncludeRoutesAll(opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, []*resource.Route, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}

	var all []*resource.ServiceRouteBinding
	var allRoutes []*resource.Route
	for {
		page, routes, pager, err := c.ListIncludeRoutes(opts)
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

// ListIncludeServiceInstances page all service route bindings the user has access to and include the
// associated service instances
func (c *ServiceRouteBindingClient) ListIncludeServiceInstances(opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, []*resource.ServiceInstance, *Pager, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}
	opts.Include = resource.ServiceRouteBindingIncludeNone

	var res resource.ServiceRouteBindingList
	err := c.client.get(path("/v3/service_route_bindings?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.ServiceInstances, pager, nil
}

// ListIncludeServiceInstancesAll retrieves all service route bindings the user has access to and include the
// associated service instances
func (c *ServiceRouteBindingClient) ListIncludeServiceInstancesAll(opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, []*resource.ServiceInstance, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}

	var all []*resource.ServiceRouteBinding
	var allSIs []*resource.ServiceInstance
	for {
		page, sis, pager, err := c.ListIncludeServiceInstances(opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allSIs = append(allSIs, sis...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allSIs, nil
}

// Update the specified attributes of the service route binding
func (c *ServiceRouteBindingClient) Update(guid string, r *resource.ServiceRouteBindingUpdate) (*resource.ServiceRouteBinding, error) {
	var srb resource.ServiceRouteBinding
	_, err := c.client.patch(path("/v3/service_route_bindings/%s", guid), r, &srb)
	if err != nil {
		return nil, err
	}
	return &srb, nil
}
