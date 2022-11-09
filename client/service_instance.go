package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
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

// CreateManaged requests a new service instance asynchronously from a broker. The result
// of this call is an error or the job GUID.
func (c *ServiceInstanceClient) CreateManaged(r *resource.ServiceInstanceCreate) (string, error) {
	var si resource.ServiceInstance
	jobGUID, err := c.client.post("/v3/service_instances", r, &si)
	if err != nil {
		return "", err
	}
	return jobGUID, nil
}

// CreateUserProvided creates a new user provided service instance. User provided service instances
// do not require interactions with service brokers.
func (c *ServiceInstanceClient) CreateUserProvided(r *resource.ServiceInstanceCreate) (*resource.ServiceInstance, error) {
	var si resource.ServiceInstance
	_, err := c.client.post("/v3/service_instances", r, &si)
	if err != nil {
		return nil, err
	}
	return &si, nil
}

// Delete the specified service instance returning the async deletion job GUID
func (c *ServiceInstanceClient) Delete(guid string) (string, error) {
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
