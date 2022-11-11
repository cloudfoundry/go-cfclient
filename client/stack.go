package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type StackClient commonClient

// StackListOptions list filters
type StackListOptions struct {
	*ListOptions

	Names Filter `qs:"names"` // list of stack names to filter by
}

// NewStackListOptions creates new options to pass to list
func NewStackListOptions() *StackListOptions {
	return &StackListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o StackListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new stack
func (c *StackClient) Create(r *resource.StackCreate) (*resource.Stack, error) {
	var stack resource.Stack
	_, err := c.client.post("/v3/stacks", r, &stack)
	if err != nil {
		return nil, err
	}
	return &stack, nil
}

// Delete the specified stack
func (c *StackClient) Delete(guid string) error {
	_, err := c.client.delete(path("/v3/stacks/%s", guid))
	return err
}

// Get the specified stack
func (c *StackClient) Get(guid string) (*resource.Stack, error) {
	var stack resource.Stack
	err := c.client.get(path("/v3/stacks/%s", guid), &stack)
	if err != nil {
		return nil, err
	}
	return &stack, nil
}

// List pages all stacks the user has access to
func (c *StackClient) List(opts *StackListOptions) ([]*resource.Stack, *Pager, error) {
	if opts == nil {
		opts = NewStackListOptions()
	}
	var res resource.StackList
	err := c.client.get(path("/v3/stacks?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all stacks the user has access to
func (c *StackClient) ListAll(opts *StackListOptions) ([]*resource.Stack, error) {
	if opts == nil {
		opts = NewStackListOptions()
	}
	return AutoPage[*StackListOptions, *resource.Stack](opts, func(opts *StackListOptions) ([]*resource.Stack, *Pager, error) {
		return c.List(opts)
	})
}

// ListAppsOnStack pages all apps using a given stack
func (c *StackClient) ListAppsOnStack(guid string, opts *StackListOptions) ([]*resource.App, *Pager, error) {
	if opts == nil {
		opts = NewStackListOptions()
	}
	var res resource.AppList
	err := c.client.get(path("/v3/stacks/%s/apps?%s", guid, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAppsOnStackAll retrieves all apps using a given stack
func (c *StackClient) ListAppsOnStackAll(guid string, opts *StackListOptions) ([]*resource.App, error) {
	if opts == nil {
		opts = NewStackListOptions()
	}
	return AutoPage[*StackListOptions, *resource.App](opts, func(opts *StackListOptions) ([]*resource.App, *Pager, error) {
		return c.ListAppsOnStack(guid, opts)
	})
}

// Update the specified attributes of a stack
func (c *StackClient) Update(guid string, r *resource.StackUpdate) (*resource.Stack, error) {
	var stack resource.Stack
	_, err := c.client.patch(path("/v3/stacks/%s", guid), r, &stack)
	if err != nil {
		return nil, err
	}
	return &stack, nil
}
