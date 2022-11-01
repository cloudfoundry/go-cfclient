package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"net/url"
)

type StackClient commonClient

// StackListOptions list filters
type StackListOptions struct {
	*ListOptions

	Names Filter `filter:"names,omitempty"` // list of stack names to filter by
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

// Create a new space
func (c *StackClient) Create(r *resource.StackCreate) (*resource.Stack, error) {
	var space resource.Stack
	err := c.client.post(r.Name, "/v3/stacks", r, &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}

// Delete the specified space
func (c *StackClient) Delete(guid string) error {
	return c.client.delete(path("/v3/stacks/%s", guid))
}

// Get the specified space
func (c *StackClient) Get(guid string) (*resource.Stack, error) {
	var space resource.Stack
	err := c.client.get(path("/v3/stacks/%s", guid), &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}

// List pages all spaces the user has access to
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

// ListAll retrieves all spaces the user has access to
func (c *StackClient) ListAll(opts *StackListOptions) ([]*resource.Stack, error) {
	if opts == nil {
		opts = NewStackListOptions()
	}
	return AutoPage[*StackListOptions, *resource.Stack](opts, func(opts *StackListOptions) ([]*resource.Stack, *Pager, error) {
		return c.List(opts)
	})
}

// Update the specified attributes of a space
func (c *StackClient) Update(guid string, r *resource.StackUpdate) (*resource.Stack, error) {
	var space resource.Stack
	err := c.client.patch(path("/v3/stacks/%s", guid), r, &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}
