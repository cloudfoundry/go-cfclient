package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type RoleClient commonClient

// RoleIncludeType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#include
type RoleIncludeType int

const (
	RoleIncludeNone RoleIncludeType = iota
	RoleIncludeUser
	RoleIncludeSpace
	RoleIncludeOrganization
)

func (r RoleIncludeType) String() string {
	switch r {
	case RoleIncludeUser:
		return "user"
	case RoleIncludeSpace:
		return "space"
	case RoleIncludeOrganization:
		return "organization"
	}
	return ""
}

func (r RoleIncludeType) ToQueryString() url.Values {
	v := url.Values{}
	if r != RoleIncludeNone {
		v.Set("include", r.String())
	}
	return v
}

// RoleListOptions list filters
type RoleListOptions struct {
	*ListOptions

	GUIDs             Filter `filter:"guids,omitempty"`              // list of role guids to filter by
	Types             Filter `filter:"types,omitempty"`              //  list of role types to filter by
	OrganizationGUIDs Filter `filter:"organization_guids,omitempty"` // list of org guids to filter by
	SpaceGUIDs        Filter `filter:"space_guids,omitempty"`        // list of space guids to filter by
	UserGUIDs         Filter `filter:"user_guids,omitempty"`         // list of user guids to filter by
}

// NewRoleListOptions creates new options to pass to list
func NewRoleListOptions() *RoleListOptions {
	return &RoleListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o RoleListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// User returns only the specified user's roles
func (o *RoleListOptions) User(userGUID string) *RoleListOptions {
	o.UserGUIDs = Filter{
		Values: []string{userGUID},
	}
	return o
}

// Space returns only the specified space's roles
func (o *RoleListOptions) Space(spaceGUID string) *RoleListOptions {
	o.SpaceGUIDs = Filter{
		Values: []string{spaceGUID},
	}
	return o
}

// Organization returns only the specified organization's roles
func (o *RoleListOptions) Organization(orgGUID string) *RoleListOptions {
	o.OrganizationGUIDs = Filter{
		Values: []string{orgGUID},
	}
	return o
}

// OrganizationRoleType returns only roles with the specified org role type
func (o *RoleListOptions) OrganizationRoleType(roleType resource.OrganizationRoleType) *RoleListOptions {
	o.Types = Filter{
		Values: []string{roleType.String()},
	}
	return o
}

// SpaceRoleType returns only roles with the specified space role type
func (o *RoleListOptions) SpaceRoleType(roleType resource.SpaceRoleType) *RoleListOptions {
	o.Types = Filter{
		Values: []string{roleType.String()},
	}
	return o
}

// RoleListIncludeOptions list filters
type RoleListIncludeOptions struct {
	*RoleListOptions

	Include RoleIncludeType `filter:"include,omitempty"`
}

// NewRoleListIncludeOptions creates new options to pass to list
func NewRoleListIncludeOptions(include RoleIncludeType) *RoleListIncludeOptions {
	return &RoleListIncludeOptions{
		Include:         include,
		RoleListOptions: NewRoleListOptions(),
	}
}

func (o RoleListIncludeOptions) ToQueryString() url.Values {
	u := o.RoleListOptions.ToQueryString()
	if o.Include != RoleIncludeNone {
		u.Set("include", o.Include.String())
	}
	return u
}

// CreateSpaceRole creates a new role for a user in the space
//
// To create a space role you must be an admin, an organization manager
// in the parent organization of the space associated with the role,
// or a space manager in the space associated with the role.
//
// For a user to be assigned a space role, the user must already
// have an organization role in the parent organization.
func (c *RoleClient) CreateSpaceRole(spaceGUID, userGUID string, roleType resource.SpaceRoleType) (*resource.Role, error) {
	req := resource.NewRoleSpaceCreate(spaceGUID, userGUID, roleType)
	var r resource.Role
	err := c.client.post(req.RoleType, "/v3/roles", req, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateOrganizationRole creates a new role for a user in the organization
//
// To create an organization role you must be an admin or organization
// manager in the organization associated with the role.
func (c *RoleClient) CreateOrganizationRole(orgGUID, userGUID string, roleType resource.OrganizationRoleType) (*resource.Role, error) {
	req := resource.NewRoleOrganizationCreate(orgGUID, userGUID, roleType)
	var r resource.Role
	err := c.client.post(req.RoleType, "/v3/roles", req, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Delete the specified role
func (c *RoleClient) Delete(guid string) error {
	return c.client.delete(path("/v3/roles/%s", guid))
}

// Get the specified role
func (c *RoleClient) Get(guid string) (*resource.Role, error) {
	var r resource.Role
	err := c.client.get(path("/v3/roles/%s", guid), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// List all roles the user has access to in paged results
func (c *RoleClient) List(opts *RoleListOptions) ([]*resource.Role, *Pager, error) {
	if opts == nil {
		opts = NewRoleListOptions()
	}
	var res resource.RoleList
	err := c.client.get(path("/v3/roles?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all roles the user has access to
func (c *RoleClient) ListAll(opts *RoleListOptions) ([]*resource.Role, error) {
	if opts == nil {
		opts = NewRoleListOptions()
	}
	return AutoPage[*RoleListOptions, *resource.Role](opts, func(opts *RoleListOptions) ([]*resource.Role, *Pager, error) {
		return c.List(opts)
	})
}

// ListInclude pages all roles and specified included parent types the user has access to
func (c *RoleClient) ListInclude(opts *RoleListIncludeOptions) ([]*resource.Role, *resource.RoleIncluded, *Pager, error) {
	if opts == nil {
		opts = NewRoleListIncludeOptions(RoleIncludeNone)
	}
	var res resource.RoleList
	err := c.client.get(path("/v3/roles?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included, pager, nil
}

// ListIncludeAll retrieves all roles and specified included parent types the user has access to
func (c *RoleClient) ListIncludeAll(opts *RoleListIncludeOptions) ([]*resource.Role, *resource.RoleIncluded, error) {
	if opts == nil {
		opts = NewRoleListIncludeOptions(RoleIncludeNone)
	}
	return roleAutoPageInclude[*RoleListIncludeOptions, *resource.Role](opts, func(opts *RoleListIncludeOptions) ([]*resource.Role, *resource.RoleIncluded, *Pager, error) {
		return c.ListInclude(opts)
	})
}

type roleListIncludeFunc[T ListOptioner, R any] func(opts T) ([]R, *resource.RoleIncluded, *Pager, error)

func roleAutoPageInclude[T ListOptioner, R any](opts T, list roleListIncludeFunc[T, R]) ([]R, *resource.RoleIncluded, error) {
	var all []R
	var allIncluded *resource.RoleIncluded
	for {
		page, included, pager, err := list(opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allIncluded.Organizations = append(allIncluded.Organizations, included.Organizations...)
		allIncluded.Spaces = append(allIncluded.Spaces, included.Spaces...)
		allIncluded.Users = append(allIncluded.Users, included.Users...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allIncluded, nil
}
