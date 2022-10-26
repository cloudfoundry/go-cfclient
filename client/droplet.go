package client

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"io"
	"net/http"
)

type DropletClient commonClient

// DropletListOptions list filters
type DropletListOptions struct {
	*ListOptions

	GUIDs             Filter `filter:"guids,omitempty"`              // list of droplet guids to filter by
	States            Filter `filter:"states,omitempty"`             // list of droplet states to filter by
	AppGUIDs          Filter `filter:"app_guids,omitempty"`          // list of app guids to filter by
	SpaceGUIDs        Filter `filter:"space_guids,omitempty"`        // list of space guids to filter by
	OrganizationGUIDs Filter `filter:"organization_guids,omitempty"` // list of organization guids to filter by
}

// NewDropletListOptions creates new options to pass to list
func NewDropletListOptions() *DropletListOptions {
	return &DropletListOptions{
		ListOptions: NewListOptions(),
	}
}

// DropletPackageListOptions list filters
type DropletPackageListOptions struct {
	*ListOptions

	GUIDs  Filter `filter:"guids,omitempty"`  // list of droplet guids to filter by
	States Filter `filter:"states,omitempty"` // list of droplet states to filter by
}

// NewDropletPackageListOptions creates new options to pass to list droplets by package
func NewDropletPackageListOptions() *DropletPackageListOptions {
	return &DropletPackageListOptions{
		ListOptions: NewListOptions(),
	}
}

// DropletAppListOptions list filters
type DropletAppListOptions struct {
	*ListOptions

	GUIDs   Filter `filter:"guids,omitempty"`   // list of droplet guids to filter by
	States  Filter `filter:"states,omitempty"`  // list of droplet states to filter by
	Current bool   `filter:"current,omitempty"` // If true, only include the droplet currently assigned to the app
}

// NewDropletAppListOptions creates new options to pass to list droplets by package
func NewDropletAppListOptions() *DropletAppListOptions {
	return &DropletAppListOptions{
		ListOptions: NewListOptions(),
	}
}

// Copy a droplet to a different app. The copied droplet excludes the environment variables listed on the source droplet
func (c *DropletClient) Copy(srcDropletGUID string, destAppGUID string) (any, error) {
	var d resource.Droplet
	r := resource.NewDropletCopy(destAppGUID)
	err := c.client.post(srcDropletGUID, path("/v3/droplets?source_guid=%s", srcDropletGUID), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Create a droplet without a package. To create a droplet based on a package, see Create a build
func (c *DropletClient) Create(r *resource.DropletCreate) (*resource.Droplet, error) {
	var d resource.Droplet
	err := c.client.post(r.Relationships.App.Data.GUID, "/v3/droplets", r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Delete the specified droplet
func (c *DropletClient) Delete(guid string) error {
	return c.client.delete(path("/v3/droplets/%s", guid))
}

// Download a gzip compressed tarball file containing a Cloud Foundry compatible droplet
// It is the caller's responsibility to close the io.ReadCloser
func (c *DropletClient) Download(guid string) (io.ReadCloser, error) {
	// This is the initial request, which will redirect to the internal blobstore location.
	// The client should automatically follow this redirect. External blob stores are untested.
	// https://v3-apidocs.cloudfoundry.org/version/3.127.0/index.html#download-droplet-bits
	p := path("/v3/droplets/%s/download", guid)
	req := c.client.NewRequest("GET", p)
	resp, err := c.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error getting %s: %w", p, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting %s, response code: %d", p, resp.StatusCode)
	}
	return resp.Body, nil
}

// Get retrieves the droplet by ID
func (c *DropletClient) Get(guid string) (*resource.Droplet, error) {
	var d resource.Droplet
	err := c.client.get(path("/v3/droplets/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// List pages all droplets the user has access to
func (c *DropletClient) List(opts *DropletListOptions) ([]*resource.Droplet, *Pager, error) {
	var res resource.DropletList
	err := c.client.get(path("/v3/droplets?%s", opts.ToQueryString(opts)), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all droplets the user has access to
func (c *DropletClient) ListAll() ([]*resource.Droplet, error) {
	opts := NewDropletListOptions()
	var allDroplets []*resource.Droplet
	for {
		apps, pager, err := c.List(opts)
		if err != nil {
			return nil, err
		}
		allDroplets = append(allDroplets, apps...)
		if !pager.HasNextPage() {
			break
		}
		opts.ListOptions = pager.NextPage(opts.ListOptions)
	}
	return allDroplets, nil
}

// ListForApp pages all droplets for the specified app
func (c *DropletClient) ListForApp(appGUID string, opts *DropletAppListOptions) ([]*resource.Droplet, *Pager, error) {
	var res resource.DropletList
	err := c.client.get(path("/v3/apps/%s/droplets?%s", appGUID, opts.ToQueryString(opts)), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForPackage pages all droplets for the specified package
func (c *DropletClient) ListForPackage(packageGUID string, opts *DropletPackageListOptions) ([]*resource.Droplet, *Pager, error) {
	var res resource.DropletList
	err := c.client.get(path("/v3/packages/%s/droplets?%s", packageGUID, opts.ToQueryString(opts)), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// GetCurrentAssociationForApp retrieves the current droplet relationship for an app
func (c *DropletClient) GetCurrentAssociationForApp(appGUID string) (*resource.CurrentDropletResponse, error) {
	var d resource.CurrentDropletResponse
	err := c.client.get(path("/v3/apps/%s/relationships/current_droplet", appGUID), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetCurrentForApp retrieves the current droplet for an app
func (c *DropletClient) GetCurrentForApp(appGUID string) (*resource.Droplet, error) {
	var d resource.Droplet
	err := c.client.get(path("/v3/apps/%s/droplets/current", appGUID), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// SetCurrentAssociationForApp sets the current droplet for an app. The current droplet is the droplet that the app will use when running
func (c *DropletClient) SetCurrentAssociationForApp(appGUID, dropletGUID string) (*resource.CurrentDropletResponse, error) {
	var d resource.CurrentDropletResponse
	r := resource.ToOneRelationship{Data: resource.Relationship{GUID: dropletGUID}}
	err := c.client.patch(path("/v3/apps/%s/relationships/current_droplet", appGUID), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Update an existing droplet
func (c *DropletClient) Update(guid string, r *resource.DropletUpdate) (*resource.Droplet, error) {
	var d resource.Droplet
	err := c.client.patch(path("/v3/droplets/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
