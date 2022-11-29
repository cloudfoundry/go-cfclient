package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type OrgClient commonClient

type OrgListOptions struct {
	*ListOptions

	GUIDs Filter `qs:"guids"` // list of organization guids to filter by
	Names Filter `qs:"names"` // list of organization names to filter by
}

// NewOrgListOptions creates new options to pass to list
func NewOrgListOptions() *OrgListOptions {
	return &OrgListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o OrgListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// AssignDefaultIsoSegment assigns a default iso segment to the specified org
//
// Apps will not run in the new default isolation segment until they are restarted
func (c *OrgClient) AssignDefaultIsoSegment(ctx context.Context, guid, isoSegmentGUID string) error {
	r := &resource.ToOneRelationship{
		Data: &resource.Relationship{
			GUID: isoSegmentGUID,
		},
	}
	_, err := c.client.patch(ctx, path.Format("/v3/organizations/%s/relationships/default_isolation_segment", guid), r, nil)
	return err
}

// Create an organization
func (c *OrgClient) Create(ctx context.Context, r *resource.OrganizationCreate) (*resource.Organization, error) {
	var org resource.Organization
	_, err := c.client.post(ctx, "/v3/organizations", r, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// Delete the specified organization asynchronously and return a jobGUID
func (c *OrgClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/organizations/%s", guid))
}

// First returns the first organization matching the options or an error when less than 1 match
func (c *OrgClient) First(ctx context.Context, opts *OrgListOptions) (*resource.Organization, error) {
	return First[*OrgListOptions, *resource.Organization](opts, func(opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// FirstForIsoSegment returns the first organization matching the options and iso segment or an error when less than 1 match
func (c *OrgClient) FirstForIsoSegment(ctx context.Context, isoSegmentGUID string, opts *OrgListOptions) (*resource.Organization, error) {
	return First[*OrgListOptions, *resource.Organization](opts, func(opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
		return c.ListForIsoSegment(ctx, isoSegmentGUID, opts)
	})
}

// Get the specified organization
func (c *OrgClient) Get(ctx context.Context, guid string) (*resource.Organization, error) {
	var org resource.Organization
	err := c.client.get(ctx, path.Format("/v3/organizations/%s", guid), &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// GetDefaultIsoSegment gets the specified organization's default iso segment GUID if any
func (c *OrgClient) GetDefaultIsoSegment(ctx context.Context, guid string) (string, error) {
	var relation resource.ToOneRelationship
	err := c.client.get(ctx, path.Format("/v3/organizations/%s/relationships/default_isolation_segment", guid), &relation)
	if err != nil {
		return "", err
	}
	return relation.Data.GUID, nil
}

// GetDefaultDomain gets the specified organization's default domain if any
func (c *OrgClient) GetDefaultDomain(ctx context.Context, guid string) (*resource.Domain, error) {
	var domain resource.Domain
	err := c.client.get(ctx, path.Format("/v3/organizations/%s/domains/default", guid), &domain)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

// GetUsageSummary gets the specified organization's usage summary
func (c *OrgClient) GetUsageSummary(ctx context.Context, guid string) (*resource.OrganizationUsageSummary, error) {
	var summary resource.OrganizationUsageSummary
	err := c.client.get(ctx, path.Format("/v3/organizations/%s/usage_summary", guid), &summary)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// List pages all organizations the user has access to
func (c *OrgClient) List(ctx context.Context, opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewOrgListOptions()
	}
	var res resource.OrganizationList
	err := c.client.get(ctx, path.Format("/v3/organizations?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all organizations the user has access to
func (c *OrgClient) ListAll(ctx context.Context, opts *OrgListOptions) ([]*resource.Organization, error) {
	if opts == nil {
		opts = NewOrgListOptions()
	}
	return AutoPage[*OrgListOptions, *resource.Organization](opts, func(opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListForIsoSegment pages all organizations for the specified isolation segment
func (c *OrgClient) ListForIsoSegment(ctx context.Context, isoSegmentGUID string, opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewOrgListOptions()
	}
	var res resource.OrganizationList
	err := c.client.get(ctx, path.Format("/v3/isolation_segments/%s/organizations?%s", isoSegmentGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForIsoSegmentAll retrieves all organizations for the specified isolation segment
func (c *OrgClient) ListForIsoSegmentAll(ctx context.Context, isoSegmentGUID string, opts *OrgListOptions) ([]*resource.Organization, error) {
	if opts == nil {
		opts = NewOrgListOptions()
	}
	return AutoPage[*OrgListOptions, *resource.Organization](opts, func(opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
		return c.ListForIsoSegment(ctx, isoSegmentGUID, opts)
	})
}

// ListUsers pages of all users that are members of the specified org
func (c *OrgClient) ListUsers(ctx context.Context, guid string, opts *UserListOptions) ([]*resource.User, *Pager, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	var res resource.UserList
	err := c.client.get(ctx, path.Format("/v3/organizations/%s/users?%s", guid, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListUsersAll retrieves all users that are members of the specified org
func (c *OrgClient) ListUsersAll(ctx context.Context, guid string, opts *UserListOptions) ([]*resource.User, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	return AutoPage[*UserListOptions, *resource.User](opts, func(opts *UserListOptions) ([]*resource.User, *Pager, error) {
		return c.ListUsers(ctx, guid, opts)
	})
}

// Single returns a single org matching the options or an error if not exactly 1 match
func (c *OrgClient) Single(ctx context.Context, opts *OrgListOptions) (*resource.Organization, error) {
	return Single[*OrgListOptions, *resource.Organization](opts, func(opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// SingleForIsoSegment returns a single org matching the options and iso segment or an error if not exactly 1 match
func (c *OrgClient) SingleForIsoSegment(ctx context.Context, isoSegmentGUID string, opts *OrgListOptions) (*resource.Organization, error) {
	return Single[*OrgListOptions, *resource.Organization](opts, func(opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
		return c.ListForIsoSegment(ctx, isoSegmentGUID, opts)
	})
}

// Update the organization's specified attributes
func (c *OrgClient) Update(ctx context.Context, guid string, r *resource.OrganizationUpdate) (*resource.Organization, error) {
	var org resource.Organization
	_, err := c.client.patch(ctx, path.Format("/v3/organizations/%s", guid), r, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}
