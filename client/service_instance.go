package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"net/url"
)

type ServiceInstanceClient commonClient

// ServiceInstanceListOptions list filters
type ServiceInstanceListOptions struct {
	*ListOptions

	Names             Filter `filter:"names,omitempty"` // list of service instance names to filter by
	GUIDs             Filter `filter:"guids,omitempty"` // list of service instance guids to filter by
	Type              string `filter:"type,omitempty"`  // Filter by type; valid values are managed and user-provided
	SpaceGUIDs        Filter `filter:"space_guids,omitempty"`
	OrganizationGUIDs Filter `filter:"organization_guids,omitempty"`
	ServicePlanGUIDs  Filter `filter:"service_plan_guids,omitempty"`
	ServicePlanNames  Filter `filter:"service_plan_names,omitempty"`
}

// NewServiceInstanceListOptions creates new options to pass to list
func NewServiceInstanceListOptions() *ServiceInstanceListOptions {
	return &ServiceInstanceListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceInstanceListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new service instance
func (c *ServiceInstanceClient) Create(r *resource.ServiceInstanceCreate) (*resource.ServiceInstance, error) {
	var si resource.ServiceInstance
	err := c.client.post(r.Name, "/v3/service_instances", r, &si)
	if err != nil {
		return nil, err
	}
	return &si, nil
}

// Delete the specified service instance
func (c *ServiceInstanceClient) Delete(guid string) error {
	return c.client.delete(path("/v3/service_instances/%s", guid))
}

// Get the specified service instance
func (c *ServiceInstanceClient) Get(guid string) (*resource.ServiceInstance, error) {
	var si resource.ServiceInstance
	err := c.client.get(path("/v3/service_instances/%s", guid), &si)
	if err != nil {
		return nil, err
	}
	return &si, nil
}

// List pages all service instances the user has access to
func (c *ServiceInstanceClient) List(opts *ServiceInstanceListOptions) ([]*resource.ServiceInstance, *Pager, error) {
	if opts == nil {
		opts = NewServiceInstanceListOptions()
	}
	var res resource.ServiceInstanceList
	err := c.client.get(path("/v3/service_instances?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service instances the user has access to
func (c *ServiceInstanceClient) ListAll(opts *ServiceInstanceListOptions) ([]*resource.ServiceInstance, error) {
	if opts == nil {
		opts = NewServiceInstanceListOptions()
	}
	return AutoPage[*ServiceInstanceListOptions, *resource.ServiceInstance](opts, func(opts *ServiceInstanceListOptions) ([]*resource.ServiceInstance, *Pager, error) {
		return c.List(opts)
	})
}
