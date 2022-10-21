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

type SpaceListOptions struct {
	*ListOptions

	GUIDs             Filter
	Names             Filter
	OrganizationGUIDs Filter
	Include           SpaceIncludeType
}

func NewSpaceListOptions() *SpaceListOptions {
	return &SpaceListOptions{
		ListOptions: NewListOptions(),
	}
}

func (c *SpaceClient) Create(r *resource.SpaceCreate) (*resource.Space, error) {
	var space resource.Space
	err := c.client.post(r.Name, "/v3/spaces", r, &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}

func (c *SpaceClient) Delete(guid string) error {
	return c.client.delete(path("/v3/spaces/%s", guid))
}

func (c *SpaceClient) Get(guid string) (*resource.Space, error) {
	var space resource.Space
	err := c.client.get(path("/v3/spaces/%s", guid), &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}

func (c *SpaceClient) GetAndInclude(guid string, include SpaceIncludeType) (*resource.Space, error) {
	var space resource.Space
	err := c.client.get(path("/v3/spaces/%s?%s", guid, include.ToQueryString()), &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}

func (c *SpaceClient) List(opts *SpaceListOptions) ([]*resource.Space, *Pager, error) {
	var res resource.SpaceList
	err := c.client.get(path("/v3/spaces?%s", opts.ToQueryString(opts)), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

func (c *SpaceClient) ListAll() ([]*resource.Space, error) {
	opts := NewSpaceListOptions()
	var allSpaces []*resource.Space
	for {
		spaces, pager, err := c.List(opts)
		if err != nil {
			return nil, err
		}
		allSpaces = append(allSpaces, spaces...)
		if !pager.HasNextPage() {
			break
		}
		opts.ListOptions = pager.NextPage(opts.ListOptions)
	}
	return allSpaces, nil
}

// ListUsers lists users by space GUID
func (c *SpaceClient) ListUsers(spaceGUID string) ([]*resource.User, *Pager, error) {
	var res resource.SpaceUserList
	err := c.client.get(path("/v3/spaces/%s/users", spaceGUID), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

func (c *SpaceClient) ListUsersAll(spaceGUID string) ([]*resource.User, error) {
	opts := NewListOptions()
	var allUsers []*resource.User
	for {
		users, pager, err := c.ListUsers(spaceGUID)
		if err != nil {
			return nil, err
		}
		allUsers = append(allUsers, users...)
		if !pager.HasNextPage() {
			break
		}
		opts = pager.NextPage(opts)
	}
	return allUsers, nil
}

func (c *SpaceClient) Update(guid string, r *resource.SpaceUpdate) (*resource.Space, error) {
	var space resource.Space
	err := c.client.patch(path("/v3/spaces/%s", guid), r, &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}
