package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type RoleClient commonClient

// RoleListOptions list filters
type RoleListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`              // list of role guids to filter by
	Types             Filter `qs:"types"`              // list of role types to filter by
	OrganizationGUIDs Filter `qs:"organization_guids"` // list of org guids to filter by
	SpaceGUIDs        Filter `qs:"space_guids"`        // list of space guids to filter by
	UserGUIDs         Filter `qs:"user_guids"`         // list of user guids to filter by

	Include resource.RoleIncludeType `qs:"include"`
}

// NewRoleListOptions creates new options to pass to list
func NewRoleListOptions() *RoleListOptions {
	return &RoleListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o *RoleListOptions) ToQueryString() url.Values {
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

// CreateSpaceRole creates a new role for a user in the space
//
// To create a space role you must be an admin, an organization manager
// in the parent organization of the space associated with the role,
// or a space manager in the space associated with the role.
//
// For a user to be assigned a space role, the user must already
// have an organization role in the parent organization.
func (c *RoleClient) CreateSpaceRole(ctx context.Context, spaceGUID, userGUID string, roleType resource.SpaceRoleType) (*resource.Role, error) {
	req := resource.NewRoleSpaceCreate(spaceGUID, userGUID, roleType)
	var r resource.Role
	_, err := c.client.post(ctx, "/v3/roles", req, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateOrganizationRole creates a new role for a user in the organization
//
// To create an organization role you must be an admin or organization
// manager in the organization associated with the role.
func (c *RoleClient) CreateOrganizationRole(ctx context.Context, orgGUID, userGUID string, roleType resource.OrganizationRoleType) (*resource.Role, error) {
	req := resource.NewRoleOrganizationCreate(orgGUID, userGUID, roleType)
	var r resource.Role
	_, err := c.client.post(ctx, "/v3/roles", req, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Delete the specified role asynchronously and return a jobGUID
func (c *RoleClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/roles/%s", guid))
}

// Get the specified role
func (c *RoleClient) Get(ctx context.Context, guid string) (*resource.Role, error) {
	var r resource.Role
	err := c.client.get(ctx, path.Format("/v3/roles/%s", guid), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// GetIncludeOrgs allows callers to fetch a role and include any assigned orgs
func (c *RoleClient) GetIncludeOrgs(ctx context.Context, guid string) (*resource.Role, []*resource.Organization, error) {
	var role resource.RoleWithIncluded
	err := c.client.get(ctx, path.Format("/v3/roles/%s?include=%s", guid, resource.RoleIncludeOrganization), &role)
	if err != nil {
		return nil, nil, err
	}
	return &role.Role, role.Included.Organizations, nil
}

// GetIncludeSpaces allows callers to fetch a role and include any assigned spaces
func (c *RoleClient) GetIncludeSpaces(ctx context.Context, guid string) (*resource.Role, []*resource.Space, error) {
	var role resource.RoleWithIncluded
	err := c.client.get(ctx, path.Format("/v3/roles/%s?include=%s", guid, resource.RoleIncludeSpace), &role)
	if err != nil {
		return nil, nil, err
	}
	return &role.Role, role.Included.Spaces, nil
}

// GetIncludeUsers allows callers to fetch a role and include any assigned users
func (c *RoleClient) GetIncludeUsers(ctx context.Context, guid string) (*resource.Role, []*resource.User, error) {
	var role resource.RoleWithIncluded
	err := c.client.get(ctx, path.Format("/v3/roles/%s?include=%s", guid, resource.RoleIncludeUser), &role)
	if err != nil {
		return nil, nil, err
	}
	return &role.Role, role.Included.Users, nil
}

// List all roles the user has access to in paged results
func (c *RoleClient) List(ctx context.Context, opts *RoleListOptions) ([]*resource.Role, *Pager, error) {
	if opts == nil {
		opts = NewRoleListOptions()
	}
	var res resource.RoleList
	err := c.client.get(ctx, path.Format("/v3/roles?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all roles the user has access to
func (c *RoleClient) ListAll(ctx context.Context, opts *RoleListOptions) ([]*resource.Role, error) {
	if opts == nil {
		opts = NewRoleListOptions()
	}
	return AutoPage[*RoleListOptions, *resource.Role](opts, func(opts *RoleListOptions) ([]*resource.Role, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListIncludeOrgs pages all roles and specified and includes orgs that have the roles
func (c *RoleClient) ListIncludeOrgs(ctx context.Context, opts *RoleListOptions) ([]*resource.Role, []*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewRoleListOptions()
	}
	opts.Include = resource.RoleIncludeOrganization

	var res resource.RoleList
	err := c.client.get(ctx, path.Format("/v3/roles?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Organizations, pager, nil
}

// ListIncludeOrgsAll retrieves all roles and specified and includes orgs that have the roles
func (c *RoleClient) ListIncludeOrgsAll(ctx context.Context, opts *RoleListOptions) ([]*resource.Role, []*resource.Organization, error) {
	if opts == nil {
		opts = NewRoleListOptions()
	}
	opts.Include = resource.RoleIncludeOrganization

	var all []*resource.Role
	var allOrgs []*resource.Organization
	for {
		page, orgs, pager, err := c.ListIncludeOrgs(ctx, opts)
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

// ListIncludeSpaces pages all roles and specified and includes spaces that have the roles
func (c *RoleClient) ListIncludeSpaces(ctx context.Context, opts *RoleListOptions) ([]*resource.Role, []*resource.Space, *Pager, error) {
	if opts == nil {
		opts = NewRoleListOptions()
	}
	opts.Include = resource.RoleIncludeSpace

	var res resource.RoleList
	err := c.client.get(ctx, path.Format("/v3/roles?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Spaces, pager, nil
}

// ListIncludeSpacesAll retrieves all roles and specified and includes spaces that have the roles
func (c *RoleClient) ListIncludeSpacesAll(ctx context.Context, opts *RoleListOptions) ([]*resource.Role, []*resource.Space, error) {
	if opts == nil {
		opts = NewRoleListOptions()
	}
	opts.Include = resource.RoleIncludeSpace

	var all []*resource.Role
	var allSpaces []*resource.Space
	for {
		page, spaces, pager, err := c.ListIncludeSpaces(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allSpaces = append(allSpaces, spaces...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allSpaces, nil
}

// ListIncludeUsers pages all roles and specified and includes users that belong to the roles
func (c *RoleClient) ListIncludeUsers(ctx context.Context, opts *RoleListOptions) ([]*resource.Role, []*resource.User, *Pager, error) {
	if opts == nil {
		opts = NewRoleListOptions()
	}
	opts.Include = resource.RoleIncludeUser

	var res resource.RoleList
	err := c.client.get(ctx, path.Format("/v3/roles?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Users, pager, nil
}

// ListIncludeUsersAll retrieves all roles and all the users that belong to those roles
func (c *RoleClient) ListIncludeUsersAll(ctx context.Context, opts *RoleListOptions) ([]*resource.Role, []*resource.User, error) {
	if opts == nil {
		opts = NewRoleListOptions()
	}
	opts.Include = resource.RoleIncludeUser

	var all []*resource.Role
	var allUsers []*resource.User
	for {
		page, users, pager, err := c.ListIncludeUsers(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allUsers = append(allUsers, users...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allUsers, nil
}
