package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type IsolationSegmentClient commonClient

// IsolationSegmentOptions list filters
type IsolationSegmentOptions struct {
	*ListOptions

	GUIDs             Filter `filter:"guids,omitempty"`
	Names             Filter `filter:"names,omitempty"`
	OrganizationGUIDs Filter `filter:"organization_guids,omitempty"`
}

// NewIsolationSegmentOptions creates new options to pass to list
func NewIsolationSegmentOptions() *IsolationSegmentOptions {
	return &IsolationSegmentOptions{
		ListOptions: NewListOptions(),
	}
}

func (o IsolationSegmentOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new isolation segment
func (c *IsolationSegmentClient) Create(r *resource.IsolationSegmentCreate) (*resource.IsolationSegment, error) {
	var iso resource.IsolationSegment
	_, err := c.client.post("/v3/isolation_segments", r, &iso)
	if err != nil {
		return nil, err
	}
	return &iso, nil
}

// Delete the specified isolation segments
//
// An isolation segment cannot be deleted if it is entitled to any organization.
func (c *IsolationSegmentClient) Delete(guid string) error {
	_, err := c.client.delete(path("/v3/isolation_segments/%s", guid))
	return err
}

// EntitleOrg entitles the specified organization for the isolation segment.
//
// In the case where the specified isolation segment is the system-wide shared segment,
// and if an organization is not already entitled for any other isolation segment, then
// the shared isolation segment automatically gets assigned as the default for that organization.
func (c *IsolationSegmentClient) EntitleOrg(guid string, orgGUID string) (*resource.IsolationSegmentRelationship, error) {
	return c.EntitleOrgs(guid, []string{orgGUID})
}

// EntitleOrgs entitles the specified organizations for the isolation segment.
//
// In the case where the specified isolation segment is the system-wide shared segment,
// and if an organization is not already entitled for any other isolation segment, then
// the shared isolation segment automatically gets assigned as the default for that organization.
func (c *IsolationSegmentClient) EntitleOrgs(guid string, orgGUIDs []string) (*resource.IsolationSegmentRelationship, error) {
	req := resource.NewToManyRelationships(orgGUIDs)
	var iso resource.IsolationSegmentRelationship
	_, err := c.client.post(path("/v3/isolation_segments/%s/relationships/organizations", guid), req, &iso)
	if err != nil {
		return nil, err
	}
	return &iso, nil
}

// Get the specified isolation segment
func (c *IsolationSegmentClient) Get(guid string) (*resource.IsolationSegment, error) {
	var iso resource.IsolationSegment
	err := c.client.get(path("/v3/isolation_segments/%s", guid), &iso)
	if err != nil {
		return nil, err
	}
	return &iso, nil
}

// List all isolation segments the user has access to in paged results
//
// For admin, this is all the isolation segments in the system. For anyone else,  this is
// the isolation segments in the allowed list for any organization to which the user belongs.
func (c *IsolationSegmentClient) List(opts *IsolationSegmentOptions) ([]*resource.IsolationSegment, *Pager, error) {
	if opts == nil {
		opts = NewIsolationSegmentOptions()
	}

	var isos resource.IsolationSegmentList
	err := c.client.get(path("/v3/isolation_segments?%s", opts.ToQueryString()), &isos)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(isos.Pagination)
	return isos.Resources, pager, nil
}

// ListAll retrieves all isolation segments the user has access to
//
// For admin, this is all the isolation segments in the system. For anyone else,  this is
// the isolation segments in the allowed list for any organization to which the user belongs.
func (c *IsolationSegmentClient) ListAll(opts *IsolationSegmentOptions) ([]*resource.IsolationSegment, error) {
	if opts == nil {
		opts = NewIsolationSegmentOptions()
	}
	return AutoPage[*IsolationSegmentOptions, *resource.IsolationSegment](opts, func(opts *IsolationSegmentOptions) ([]*resource.IsolationSegment, *Pager, error) {
		return c.List(opts)
	})
}

// ListOrgRelationships lists the organizations entitled for the isolation segment.
//
// For an Admin, this will list all entitled organizations in the system. For any other user,
// this will list only the entitled organizations to which the user belongs.
func (c *IsolationSegmentClient) ListOrgRelationships(guid string) ([]string, error) {
	var relationships resource.IsolationSegmentRelationship
	err := c.client.get(path("/v3/isolation_segments/%s/relationships/organizations", guid), &relationships)
	if err != nil {
		return nil, err
	}

	var orgGUIDs []string
	for _, relation := range relationships.Data {
		orgGUIDs = append(orgGUIDs, relation.GUID)
	}
	return orgGUIDs, nil
}

// ListSpaceRelationships lists the spaces to which the isolation segment is assigned.
//
// For an Admin, this will list all associated spaces in the system. For an org manager,
// this will list only those associated spaces belonging to orgs for which the user is a
// manager. For any other user, this will list only those associated spaces to which the
// user has access.
func (c *IsolationSegmentClient) ListSpaceRelationships(guid string) ([]string, error) {
	var relationships resource.IsolationSegmentRelationship
	err := c.client.get(path("/v3/isolation_segments/%s/relationships/spaces", guid), &relationships)
	if err != nil {
		return nil, err
	}

	var spaceGUIDs []string
	for _, relation := range relationships.Data {
		spaceGUIDs = append(spaceGUIDs, relation.GUID)
	}
	return spaceGUIDs, nil
}

// RevokeOrg revokes the entitlement for the specified organization to the isolation segment
//
// If the isolation segment is assigned to a space within an organization, the entitlement cannot be revoked.
// If the isolation segment is the organization’s default, the entitlement cannot be revoked.
func (c *IsolationSegmentClient) RevokeOrg(guid string, orgGUID string) error {
	_, err := c.client.delete(path("/v3/isolation_segments/%s/relationships/organizations/%s", guid, orgGUID))
	return err
}

// RevokeOrgs revokes the entitlement for all the specified organizations to the isolation segment
//
// If the isolation segment is assigned to a space within an organization, the entitlement cannot be revoked.
// If the isolation segment is the organization’s default, the entitlement cannot be revoked.
func (c *IsolationSegmentClient) RevokeOrgs(guid string, orgGUIDs []string) error {
	for _, orgGUID := range orgGUIDs {
		err := c.RevokeOrg(guid, orgGUID)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update the specified attributes of the isolation segments
func (c *IsolationSegmentClient) Update(guid string, r *resource.IsolationSegmentUpdate) (*resource.IsolationSegment, error) {
	var iso resource.IsolationSegment
	_, err := c.client.patch(path("/v3/isolation_segments/%s", guid), r, &iso)
	if err != nil {
		return nil, err
	}
	return &iso, nil
}
