package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"net/url"
)

type ServiceCredentialBindingClient commonClient

// ServiceCredentialBindingListOptions list filters
type ServiceCredentialBindingListOptions struct {
	*ListOptions

	Names                Filter `filter:"names,omitempty"`                  // list of service credential binding names to filter by
	ServiceInstanceGUIDs Filter `filter:"service_instance_guids,omitempty"` // list of SI guids to filter by
	ServiceInstanceNames Filter `filter:"service_instance_names,omitempty"` // list of SI names to filter by
	AppGUIDs             Filter `filter:"app_guids,omitempty"`              // list of app guids to filter by
	AppNames             Filter `filter:"app_names,omitempty"`              // list of app names to filter by
	ServicePlanGUIDs     Filter `filter:"service_plan_guids,omitempty"`     // list of service plan guids to filter by
	ServicePlanNames     Filter `filter:"service_plan_names,omitempty"`     // list of service plan names to filter by
	ServiceOfferingGUIDs Filter `filter:"service_offering_guids,omitempty"` // list of service offering guids to filter by
	ServiceOfferingNames Filter `filter:"service_offering_names,omitempty"` // list of service offering names to filter by
	Type                 Filter `filter:"type,omitempty"`                   // list of service credential binding types to filter by, app or key
	GUIDs                Filter `filter:"guids,omitempty"`                  // list of service route binding guids to filter by
}

// NewServiceCredentialBindingListOptions creates new options to pass to list
func NewServiceCredentialBindingListOptions() *ServiceCredentialBindingListOptions {
	return &ServiceCredentialBindingListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceCredentialBindingListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// ServiceCredentialBindingListIncludeOptions list filters
type ServiceCredentialBindingListIncludeOptions struct {
	*ServiceCredentialBindingListOptions

	Include resource.ServiceCredentialBindingIncludeType `filter:"include,omitempty"`
}

// NewServiceCredentialBindingListIncludeOptions creates new options to pass to list
func NewServiceCredentialBindingListIncludeOptions(include resource.ServiceCredentialBindingIncludeType) *ServiceCredentialBindingListIncludeOptions {
	return &ServiceCredentialBindingListIncludeOptions{
		Include:                             include,
		ServiceCredentialBindingListOptions: NewServiceCredentialBindingListOptions(),
	}
}

func (o ServiceCredentialBindingListIncludeOptions) ToQueryString() url.Values {
	u := o.ServiceCredentialBindingListOptions.ToQueryString()
	if o.Include != resource.ServiceCredentialBindingIncludeNone {
		u.Set("include", o.Include.String())
	}
	return u
}

// Create a new service credential binding
func (c *ServiceCredentialBindingClient) Create(r *resource.ServiceCredentialBindingCreate) (*resource.ServiceCredentialBinding, error) {
	var d resource.ServiceCredentialBinding
	err := c.client.post("", "/v3/service_credential_bindings", r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Delete the specified service credential binding
func (c *ServiceCredentialBindingClient) Delete(guid string) error {
	return c.client.delete(path("/v3/service_credential_bindings/%s", guid))
}

// Get the specified service credential binding
func (c *ServiceCredentialBindingClient) Get(guid string) (*resource.ServiceCredentialBinding, error) {
	var d resource.ServiceCredentialBinding
	err := c.client.get(path("/v3/service_credential_bindings/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetInclude allows callers to fetch a service credential binding and include information of parent objects in the response
func (c *ServiceCredentialBindingClient) GetInclude(guid string, include resource.ServiceCredentialBindingIncludeType) (*resource.ServiceCredentialBinding, *resource.ServiceCredentialBindingIncluded, error) {
	var r resource.ServiceCredentialBindingWithIncluded
	err := c.client.get(path("/v3/service_credential_bindings/%s?include=%s", guid, include), &r)
	if err != nil {
		return nil, nil, err
	}
	return &r.ServiceCredentialBinding, r.Included, nil
}

// List pages ServiceCredentialBindings the user has access to
func (c *ServiceCredentialBindingClient) List(opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, *Pager, error) {
	var res resource.ServiceCredentialBindingList
	err := c.client.get(path("/v3/service_credential_bindings?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all ServiceCredentialBindings the user has access to
func (c *ServiceCredentialBindingClient) ListAll(opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, error) {
	if opts == nil {
		opts = NewServiceCredentialBindingListOptions()
	}
	return AutoPage[*ServiceCredentialBindingListOptions, *resource.ServiceCredentialBinding](opts, func(opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, *Pager, error) {
		return c.List(opts)
	})
}

// ListInclude page all service credential bindings the user has access to and include the specified parent resources
func (c *ServiceCredentialBindingClient) ListInclude(opts *ServiceCredentialBindingListIncludeOptions) ([]*resource.ServiceCredentialBinding, *resource.ServiceCredentialBindingIncluded, *Pager, error) {
	if opts == nil {
		opts = NewServiceCredentialBindingListIncludeOptions(resource.ServiceCredentialBindingIncludeNone)
	}
	var res resource.ServiceCredentialBindingList
	err := c.client.get(path("/v3/service_credential_bindings?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included, pager, nil
}

// ListIncludeAll retrieves all service credential bindings the user has access to and include the specified parent resources
func (c *ServiceCredentialBindingClient) ListIncludeAll(opts *ServiceCredentialBindingListIncludeOptions) ([]*resource.ServiceCredentialBinding, *resource.ServiceCredentialBindingIncluded, error) {
	if opts == nil {
		opts = NewServiceCredentialBindingListIncludeOptions(resource.ServiceCredentialBindingIncludeNone)
	}
	return serviceCredentialBindingAutoPageInclude[*ServiceCredentialBindingListIncludeOptions, *resource.ServiceCredentialBinding](opts,
		func(opts *ServiceCredentialBindingListIncludeOptions) ([]*resource.ServiceCredentialBinding, *resource.ServiceCredentialBindingIncluded, *Pager, error) {
			return c.ListInclude(opts)
		})
}

// Update the specified attributes of the app
func (c *ServiceCredentialBindingClient) Update(guid string, r *resource.ServiceCredentialBindingUpdate) (*resource.ServiceCredentialBinding, error) {
	var d resource.ServiceCredentialBinding
	err := c.client.patch(path("/v3/service_credential_bindings/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

type serviceCredentialBindingListIncludeFunc[T ListOptioner, R any] func(opts T) ([]R, *resource.ServiceCredentialBindingIncluded, *Pager, error)

func serviceCredentialBindingAutoPageInclude[T ListOptioner, R any](opts T, list serviceCredentialBindingListIncludeFunc[T, R]) ([]R, *resource.ServiceCredentialBindingIncluded, error) {
	var all []R
	var allIncluded *resource.ServiceCredentialBindingIncluded
	for {
		page, included, pager, err := list(opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allIncluded.ServiceInstances = append(allIncluded.ServiceInstances, included.ServiceInstances...)
		allIncluded.Apps = append(allIncluded.Apps, included.Apps...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allIncluded, nil
}
