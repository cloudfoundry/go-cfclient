package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type SecurityGroupClient commonClient

// SecurityGroupListOptions list filters
type SecurityGroupListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`               // list of security group guids to filter by
	Names             Filter `qs:"names"`               // list of security group names to filter by
	RunningSpaceGUIDs Filter `qs:"running_space_guids"` // list of space guids to filter by
	StagingSpaceGUIDs Filter `qs:"staging_space_guids"` // list of space guids to filter by

	GloballyEnabledRunning *bool `qs:"globally_enabled_running"` // If true, only include the security groups that are enabled for running
	GloballyEnabledStaging *bool `qs:"globally_enabled_staging"` // If true, only include the security groups that are enabled for staging
}

// NewSecurityGroupListOptions creates new options to pass to list
func NewSecurityGroupListOptions() *SecurityGroupListOptions {
	return &SecurityGroupListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o SecurityGroupListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new domain
func (c *SecurityGroupClient) Create(ctx context.Context, r *resource.SecurityGroupCreate) (*resource.SecurityGroup, error) {
	var d resource.SecurityGroup
	_, err := c.client.post(ctx, "/v3/security_groups", r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Delete the specified security group asynchronously and return a jobGUID
func (c *SecurityGroupClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/security_groups/%s", guid))
}

// Get the specified security group
func (c *SecurityGroupClient) Get(ctx context.Context, guid string) (*resource.SecurityGroup, error) {
	var d resource.SecurityGroup
	err := c.client.get(ctx, path.Format("/v3/security_groups/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// List pages SecurityGroups the user has access to
func (c *SecurityGroupClient) List(ctx context.Context, opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, *Pager, error) {
	var res resource.SecurityGroupList
	err := c.client.get(ctx, path.Format("/v3/security_groups?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all SecurityGroups the user has access to
func (c *SecurityGroupClient) ListAll(ctx context.Context, opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, error) {
	if opts == nil {
		opts = NewSecurityGroupListOptions()
	}
	return AutoPage[*SecurityGroupListOptions, *resource.SecurityGroup](opts, func(opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Single returns a single security group matching the options or an error if not exactly 1 match
func (c *SecurityGroupClient) Single(ctx context.Context, opts *SecurityGroupListOptions) (*resource.SecurityGroup, error) {
	return Single[*SecurityGroupListOptions, *resource.SecurityGroup](opts, func(opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the app
func (c *SecurityGroupClient) Update(ctx context.Context, guid string, r *resource.SecurityGroupUpdate) (*resource.SecurityGroup, error) {
	var d resource.SecurityGroup
	_, err := c.client.patch(ctx, path.Format("/v3/security_groups/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
