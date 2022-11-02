package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type SpaceClient commonClient

// SpaceListOptions list filters
type SpaceListOptions struct {
	*ListOptions

	GUIDs             Filter `filter:"guids,omitempty"`              // list of space guids to filter by
	Names             Filter `filter:"names,omitempty"`              // list of space names to filter by
	OrganizationGUIDs Filter `filter:"organization_guids,omitempty"` // list of organization guids to filter by

	Include resource.SpaceIncludeType // include parent objects if any
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

// GetIncludeOrg allows callers to fetch a space and include the parent org
func (c *SpaceClient) GetIncludeOrg(guid string) (*resource.Space, *resource.Organization, error) {
	var space resource.SpaceWithIncluded
	err := c.client.get(path("/v3/spaces/%s?include=%s", guid, resource.SpaceIncludeOrganization), &space)
	if err != nil {
		return nil, nil, err
	}
	return &space.Space, space.Included.Organizations[0], nil
}

// List pages all spaces the user has access to
func (c *SpaceClient) List(opts *SpaceListOptions) ([]*resource.Space, *Pager, error) {
	if opts == nil {
		opts = NewSpaceListOptions()
	}
	opts.Include = resource.SpaceIncludeNone

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

// ListIncludeOrgs page all spaces the user has access to and include the parent orgs
func (c *SpaceClient) ListIncludeOrgs(opts *SpaceListOptions) ([]*resource.Space, []*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewSpaceListOptions()
	}
	opts.Include = resource.SpaceIncludeOrganization

	var res resource.SpaceList
	err := c.client.get(path("/v3/spaces?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Organizations, pager, nil
}

// ListIncludeOrgsAll retrieves all spaces the user has access to and include the parent orgs
func (c *SpaceClient) ListIncludeOrgsAll(opts *SpaceListOptions) ([]*resource.Space, []*resource.Organization, error) {
	if opts == nil {
		opts = NewSpaceListOptions()
	}

	var all []*resource.Space
	var allOrgs []*resource.Organization
	for {
		page, orgs, pager, err := c.ListIncludeOrgs(opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allOrgs = append(allOrgs, orgs...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allOrgs, nil
}

// ListUsers pages users by space GUID
func (c *SpaceClient) ListUsers(spaceGUID string, opts *UserListOptions) ([]*resource.User, *Pager, error) {
	if opts == nil {
		opts = NewUserListOptions()
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
func (c *SpaceClient) ListUsersAll(spaceGUID string, opts *UserListOptions) ([]*resource.User, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	return AutoPage[*UserListOptions, *resource.User](opts, func(opts *UserListOptions) ([]*resource.User, *Pager, error) {
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
