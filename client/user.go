package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type UserClient commonClient

// UserListOptions list filters
type UserListOptions struct {
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

// NewUserListOptions creates new options to pass to list
func NewUserListOptions() *UserListOptions {
	return &UserListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o UserListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new user
func (c *UserClient) Create(r *resource.UserCreate) (*resource.User, error) {
	var user resource.User
	_, err := c.client.post("/v3/users", r, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Delete the specified user
func (c *UserClient) Delete(guid string) error {
	_, err := c.client.delete(path("/v3/users/%s", guid))
	return err
}

// Get the specified user
func (c *UserClient) Get(guid string) (*resource.User, error) {
	var user resource.User
	err := c.client.get(path("/v3/users/%s", guid), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List pages all users the user has access to
func (c *UserClient) List(opts *UserListOptions) ([]*resource.User, *Pager, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	var res resource.UserList
	err := c.client.get(path("/v3/users?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all users the user has access to
func (c *UserClient) ListAll(opts *UserListOptions) ([]*resource.User, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	return AutoPage[*UserListOptions, *resource.User](opts, func(opts *UserListOptions) ([]*resource.User, *Pager, error) {
		return c.List(opts)
	})
}

// Update the specified attributes of a user
func (c *UserClient) Update(guid string, r *resource.UserUpdate) (*resource.User, error) {
	var user resource.User
	_, err := c.client.patch(path("/v3/users/%s", guid), r, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
