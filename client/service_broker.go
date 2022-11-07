package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type ServiceBrokerClient commonClient

// ServiceBrokerListOptions list filters
type ServiceBrokerListOptions struct {
	*ListOptions

	SpaceGUIDs Filter `filter:"space_guids,omitempty"`
	Names      Filter `filter:"names,omitempty"`
}

// NewServiceBrokerListOptions creates new options to pass to list
func NewServiceBrokerListOptions() *ServiceBrokerListOptions {
	return &ServiceBrokerListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceBrokerListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new service broker
//
// Results:
// string - the async Job GUID
func (c *ServiceBrokerClient) Create(r *resource.ServiceBrokerCreate) (string, error) {
	jobGUID, err := c.client.post("/v3/service_brokers", r, nil)
	if err != nil {
		return "", err
	}
	return jobGUID, nil
}

// Delete the specified service broker
func (c *ServiceBrokerClient) Delete(guid string) (string, error) {
	return c.client.delete(path("/v3/service_brokers/%s", guid))
}

// Get the specified service broker
func (c *ServiceBrokerClient) Get(guid string) (*resource.ServiceBroker, error) {
	var sb resource.ServiceBroker
	err := c.client.get(path("/v3/service_brokers/%s", guid), &sb)
	if err != nil {
		return nil, err
	}
	return &sb, nil
}

// List pages all the service brokers the user has access to
func (c *ServiceBrokerClient) List(opts *ServiceBrokerListOptions) ([]*resource.ServiceBroker, *Pager, error) {
	if opts == nil {
		opts = NewServiceBrokerListOptions()
	}

	var res resource.ServiceBrokerList
	err := c.client.get(path("/v3/service_brokers?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service_brokers the user has access to
func (c *ServiceBrokerClient) ListAll(opts *ServiceBrokerListOptions) ([]*resource.ServiceBroker, error) {
	if opts == nil {
		opts = NewServiceBrokerListOptions()
	}
	return AutoPage[*ServiceBrokerListOptions, *resource.ServiceBroker](opts, func(opts *ServiceBrokerListOptions) ([]*resource.ServiceBroker, *Pager, error) {
		return c.List(opts)
	})
}

// Update the specified attributes of the service_broker
func (c *ServiceBrokerClient) Update(guid string, r *resource.ServiceBrokerUpdate) (string, error) {
	var sb resource.ServiceBroker
	jobGUID, err := c.client.patch(path("/v3/service_brokers/%s", guid), r, &sb)
	if err != nil {
		return "", err
	}
	return jobGUID, nil
}
