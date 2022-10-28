package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type SpaceClient commonClient

type SpaceIncludeType int

const (
	SpaceIncludeNone SpaceIncludeType = iota
	SpaceIncludeOrganization
)

func (s SpaceIncludeType) String() string {
	switch s {
	case SpaceIncludeOrganization:
		return "organization"
	}
	return ""
}

func (s SpaceIncludeType) ToQueryString() url.Values {
	v := url.Values{}
	if s != SpaceIncludeNone {
		v.Set("include", s.String())
	}
	return v
}

// SpaceListOptions list filters
type SpaceListOptions struct {
	*ListOptions

	GUIDs             Filter           // list of space guids to filter by
	Names             Filter           // list of space names to filter by
	OrganizationGUIDs Filter           // list of organization guids to filter by
	Include           SpaceIncludeType // include parent objects if any
}

// NewSpaceListOptions creates new options to pass to list
func NewSpaceListOptions() *SpaceListOptions {
	return &SpaceListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o SpaceListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new space
func (c *SpaceClient) Create(r *resource.SpaceCreate) (*resource.Space, error) {
	var space resource.Space
	err := c.client.post(r.Name, "/v3/spaces", r, &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}

// Delete the specified space
func (c *SpaceClient) Delete(guid string) error {
	return c.client.delete(path("/v3/spaces/%s", guid))
}

// Get the specified space
func (c *SpaceClient) Get(guid string) (*resource.Space, error) {
	var space resource.Space
	err := c.client.get(path("/v3/spaces/%s", guid), &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}

// GetAndInclude allows callers to fetch an space and include information of parent objects in the response
func (c *SpaceClient) GetAndInclude(guid string, include SpaceIncludeType) (*resource.Space, error) {
	var space resource.Space
	err := c.client.get(path("/v3/spaces/%s?%s", guid, include.ToQueryString()), &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}

// List pages all spaces the user has access to
func (c *SpaceClient) List(opts *SpaceListOptions) ([]*resource.Space, *Pager, error) {
	if opts == nil {
		opts = NewSpaceListOptions()
	}
	var res resource.SpaceList
	err := c.client.get(path("/v3/spaces?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all spaces the user has access to
func (c *SpaceClient) ListAll(opts *SpaceListOptions) ([]*resource.Space, error) {
	if opts == nil {
		opts = NewSpaceListOptions()
	}
	return AutoPage[*SpaceListOptions, *resource.Space](opts, func(opts *SpaceListOptions) ([]*resource.Space, *Pager, error) {
		return c.List(opts)
	})
}

// ListUsers pages users by space GUID
func (c *SpaceClient) ListUsers(spaceGUID string) ([]*resource.User, *Pager, error) {
	var res resource.SpaceUserList
	err := c.client.get(path("/v3/spaces/%s/users", spaceGUID), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListUsersAll retrieves all users by space GUID
func (c *SpaceClient) ListUsersAll(spaceGUID string) ([]*resource.User, error) {
	opts := NewSpaceListOptions()
	return AutoPage[*SpaceListOptions, *resource.User](opts, func(opts *SpaceListOptions) ([]*resource.User, *Pager, error) {
		return c.ListUsers(spaceGUID)
	})
}

// Update the specified attributes of a space
func (c *SpaceClient) Update(guid string, r *resource.SpaceUpdate) (*resource.Space, error) {
	var space resource.Space
	err := c.client.patch(path("/v3/spaces/%s", guid), r, &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}
