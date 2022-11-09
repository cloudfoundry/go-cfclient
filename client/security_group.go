package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type SecurityGroupClient commonClient

// SecurityGroupListOptions list filters
type SecurityGroupListOptions struct {
	*ListOptions

	GUIDs             Filter `filter:"guids,omitempty"`               // list of security group guids to filter by
	Names             Filter `filter:"names,omitempty"`               // list of security group names to filter by
	RunningSpaceGUIDs Filter `filter:"running_space_guids,omitempty"` // list of space guids to filter by
	StagingSpaceGUIDs Filter `filter:"staging_space_guids,omitempty"` // list of space guids to filter by

	GloballyEnabledRunning *bool `filter:"globally_enabled_running,omitempty"` // If true, only include the security groups that are enabled for running
	GloballyEnabledStaging *bool `filter:"globally_enabled_staging,omitempty"` // If true, only include the security groups that are enabled for staging
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
func (c *SecurityGroupClient) Create(r *resource.SecurityGroupCreate) (*resource.SecurityGroup, error) {
	var d resource.SecurityGroup
	_, err := c.client.post("/v3/security_groups", r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Delete the specified security group
func (c *SecurityGroupClient) Delete(guid string) error {
	_, err := c.client.delete(path("/v3/security_groups/%s", guid))
	return err
}

// Get the specified security group
func (c *SecurityGroupClient) Get(guid string) (*resource.SecurityGroup, error) {
	var d resource.SecurityGroup
	err := c.client.get(path("/v3/security_groups/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// List pages SecurityGroups the user has access to
func (c *SecurityGroupClient) List(opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, *Pager, error) {
	var res resource.SecurityGroupList
	err := c.client.get(path("/v3/security_groups?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all SecurityGroups the user has access to
func (c *SecurityGroupClient) ListAll(opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, error) {
	if opts == nil {
		opts = NewSecurityGroupListOptions()
	}
	return AutoPage[*SecurityGroupListOptions, *resource.SecurityGroup](opts, func(opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, *Pager, error) {
		return c.List(opts)
	})
}

// Update the specified attributes of the app
func (c *SecurityGroupClient) Update(guid string, r *resource.SecurityGroupUpdate) (*resource.SecurityGroup, error) {
	var d resource.SecurityGroup
	_, err := c.client.patch(path("/v3/security_groups/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
