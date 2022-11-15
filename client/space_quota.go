package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type SpaceQuotaClient commonClient

// SpaceQuotaListOptions list filters
type SpaceQuotaListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`
	Names             Filter `qs:"names"`
	OrganizationGUIDs Filter `qs:"organization_guids"`
	SpaceGUIDs        Filter `qs:"space_guids"`
}

// NewSpaceQuotaListOptions creates new options to pass to list
func NewSpaceQuotaListOptions() *SpaceQuotaListOptions {
	return &SpaceQuotaListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o SpaceQuotaListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Apply the quota to the specified spaces
func (c *SpaceQuotaClient) Apply(guid string, spaceGUIDs []string) ([]string, error) {
	req := resource.NewToManyRelationships(spaceGUIDs)
	var relation resource.ToManyRelationships
	_, err := c.client.post(path.Format("/v3/space_quotas/%s/relationships/spaces", guid), req, &relation)
	if err != nil {
		return nil, err
	}
	var guids []string
	for _, r := range relation.Data {
		guids = append(guids, r.GUID)
	}
	return guids, nil
}

// Create a new space quota
func (c *SpaceQuotaClient) Create(r *resource.SpaceQuotaCreateOrUpdate) (*resource.SpaceQuota, error) {
	var q resource.SpaceQuota
	_, err := c.client.post("/v3/space_quotas", r, &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// Delete the specified space quota
func (c *SpaceQuotaClient) Delete(guid string) error {
	_, err := c.client.delete(path.Format("/v3/space_quotas/%s", guid))
	return err
}

// Get the specified space quota
func (c *SpaceQuotaClient) Get(guid string) (*resource.SpaceQuota, error) {
	var q resource.SpaceQuota
	err := c.client.get(path.Format("/v3/space_quotas/%s", guid), &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// List pages all space quotas the user has access to
func (c *SpaceQuotaClient) List(opts *SpaceQuotaListOptions) ([]*resource.SpaceQuota, *Pager, error) {
	if opts == nil {
		opts = NewSpaceQuotaListOptions()
	}

	var res resource.SpaceQuotaList
	err := c.client.get(path.Format("/v3/space_quotas?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all space quotas the user has access to
func (c *SpaceQuotaClient) ListAll(opts *SpaceQuotaListOptions) ([]*resource.SpaceQuota, error) {
	if opts == nil {
		opts = NewSpaceQuotaListOptions()
	}
	return AutoPage[*SpaceQuotaListOptions, *resource.SpaceQuota](opts, func(opts *SpaceQuotaListOptions) ([]*resource.SpaceQuota, *Pager, error) {
		return c.List(opts)
	})
}

// Remove the space quota from the specified space
func (c *SpaceQuotaClient) Remove(guid, spaceGUID string) error {
	_, err := c.client.delete(path.Format("/v3/space_quotas/%s/relationships/spaces/%s", guid, spaceGUID))
	return err
}

// Update the specified attributes of the org quota
func (c *SpaceQuotaClient) Update(guid string, r *resource.SpaceQuotaCreateOrUpdate) (*resource.SpaceQuota, error) {
	var q resource.SpaceQuota
	_, err := c.client.patch(path.Format("/v3/space_quotas/%s", guid), r, &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
