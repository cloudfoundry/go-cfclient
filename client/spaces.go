package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type SpaceClient commonClient

// SpaceListOptions list filters
type SpaceListOptions struct {
	*ListOptions

	GUIDs             Filter // list of space guids to filter by
	Names             Filter // list of space names to filter by
	OrganizationGUIDs Filter // list of organization guids to filter by
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

// SpaceListIncludeOptions list filters
type SpaceListIncludeOptions struct {
	*SpaceListOptions

	Include resource.SpaceIncludeType // include parent objects if any
}

// NewSpaceListIncludeOptions creates new options to pass to list
func NewSpaceListIncludeOptions(include resource.SpaceIncludeType) *SpaceListIncludeOptions {
	return &SpaceListIncludeOptions{
		Include:          include,
		SpaceListOptions: NewSpaceListOptions(),
	}
}

func (o SpaceListIncludeOptions) ToQueryString() url.Values {
	u := o.SpaceListOptions.ToQueryString()
	if o.Include != resource.SpaceIncludeNone {
		u.Set("include", o.Include.String())
	}
	return u
}

// SpaceUserListOptions list filters
type SpaceUserListOptions struct {
	*ListOptions

	// list of user guids to filter by
	GUIDs Filter `filter:"guids,omitempty"`

	// list of usernames to filter by. Mutually exclusive with partial_usernames
	UserNames Filter `filter:"usernames,omitempty"`

	// list of strings to search by. When using this query parameter, all the users that
	// contain the string provided in their username will be returned. Mutually exclusive with usernames
	PartialUsernames Filter `filter:"partial_usernames,omitempty"`

	// list of user origins (user stores) to filter by, for example, users authenticated by
	// UAA have the origin “uaa”; users authenticated by an LDAP provider have the
	// origin ldap when filtering by origins, usernames must be included
	Origins Filter `filter:"origins,omitempty"`
}

// NewSpaceUserListOptions creates new options to pass to list
func NewSpaceUserListOptions() *SpaceUserListOptions {
	return &SpaceUserListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o SpaceUserListOptions) ToQueryString() url.Values {
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

// GetInclude allows callers to fetch an space and include information of parent objects in the response
func (c *SpaceClient) GetInclude(guid string, include resource.SpaceIncludeType) (*resource.Space, *resource.SpaceIncluded, error) {
	var space resource.SpaceWithIncluded
	err := c.client.get(path("/v3/spaces/%s?include=%s", guid, include), &space)
	if err != nil {
		return nil, nil, err
	}
	return &space.Space, space.Included, nil
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

// ListInclude page all spaces the user has access to and include the specified parent resources
func (c *SpaceClient) ListInclude(opts *SpaceListIncludeOptions) ([]*resource.Space, *resource.SpaceIncluded, *Pager, error) {
	if opts == nil {
		opts = NewSpaceListIncludeOptions(resource.SpaceIncludeNone)
	}
	var res resource.SpaceList
	err := c.client.get(path("/v3/spaces?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included, pager, nil
}

// ListIncludeAll retrieves all spaces the user has access to and include the specified parent resources
func (c *SpaceClient) ListIncludeAll(opts *SpaceListIncludeOptions) ([]*resource.Space, *resource.SpaceIncluded, error) {
	if opts == nil {
		opts = NewSpaceListIncludeOptions(resource.SpaceIncludeNone)
	}
	return spaceAutoPageInclude[*SpaceListIncludeOptions, *resource.Space](opts, func(opts *SpaceListIncludeOptions) ([]*resource.Space, *resource.SpaceIncluded, *Pager, error) {
		return c.ListInclude(opts)
	})
}

// ListUsers pages users by space GUID
func (c *SpaceClient) ListUsers(spaceGUID string, opts *SpaceUserListOptions) ([]*resource.User, *Pager, error) {
	if opts == nil {
		opts = NewSpaceUserListOptions()
	}
	var res resource.SpaceUserList
	err := c.client.get(path("/v3/spaces/%s/users?%s", spaceGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListUsersAll retrieves all users by space GUID
func (c *SpaceClient) ListUsersAll(spaceGUID string, opts *SpaceUserListOptions) ([]*resource.User, error) {
	if opts == nil {
		opts = NewSpaceUserListOptions()
	}
	return AutoPage[*SpaceUserListOptions, *resource.User](opts, func(opts *SpaceUserListOptions) ([]*resource.User, *Pager, error) {
		return c.ListUsers(spaceGUID, opts)
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

type spaceListIncludeFunc[T ListOptioner, R any] func(opts T) ([]R, *resource.SpaceIncluded, *Pager, error)

func spaceAutoPageInclude[T ListOptioner, R any](opts T, list spaceListIncludeFunc[T, R]) ([]R, *resource.SpaceIncluded, error) {
	var all []R
	var allIncluded *resource.SpaceIncluded
	for {
		page, included, pager, err := list(opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allIncluded.Organizations = append(allIncluded.Organizations, included.Organizations...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allIncluded, nil
}
