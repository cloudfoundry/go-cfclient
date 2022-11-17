package client

import (
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"io"
	"net/url"
)

type DropletClient commonClient

// DropletListOptions list filters
type DropletListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`              // list of droplet guids to filter by
	States            Filter `qs:"states"`             // list of droplet states to filter by
	AppGUIDs          Filter `qs:"app_guids"`          // list of app guids to filter by
	SpaceGUIDs        Filter `qs:"space_guids"`        // list of space guids to filter by
	OrganizationGUIDs Filter `qs:"organization_guids"` // list of organization guids to filter by
}

// NewDropletListOptions creates new options to pass to list
func NewDropletListOptions() *DropletListOptions {
	return &DropletListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o DropletListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// DropletPackageListOptions list filters
type DropletPackageListOptions struct {
	*ListOptions

	GUIDs  Filter `qs:"guids"`  // list of droplet guids to filter by
	States Filter `qs:"states"` // list of droplet states to filter by
}

// NewDropletPackageListOptions creates new options to pass to list droplets by package
func NewDropletPackageListOptions() *DropletPackageListOptions {
	return &DropletPackageListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o DropletPackageListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// DropletAppListOptions list filters
type DropletAppListOptions struct {
	*ListOptions

	GUIDs   Filter `qs:"guids"`   // list of droplet guids to filter by
	States  Filter `qs:"states"`  // list of droplet states to filter by
	Current bool   `qs:"current"` // If true, only include the droplet currently assigned to the app
}

// NewDropletAppListOptions creates new options to pass to list droplets by package
func NewDropletAppListOptions() *DropletAppListOptions {
	return &DropletAppListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o DropletAppListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Copy a droplet to a different app. The copied droplet excludes the environment variables listed on the source droplet
func (c *DropletClient) Copy(srcDropletGUID string, destAppGUID string) (any, error) {
	var d resource.Droplet
	r := resource.NewDropletCopy(destAppGUID)
	_, err := c.client.post(path.Format("/v3/droplets?source_guid=%s", srcDropletGUID), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Create a droplet without a package. To create a droplet based on a package, see Create a build
func (c *DropletClient) Create(r *resource.DropletCreate) (*resource.Droplet, error) {
	var d resource.Droplet
	_, err := c.client.post("/v3/droplets", r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Delete the specified droplet asynchronously and return a jobGUID.
func (c *DropletClient) Delete(guid string) (string, error) {
	return c.client.delete(path.Format("/v3/droplets/%s", guid))
}

// Download a gzip compressed tarball file containing a Cloud Foundry compatible droplet
// It is the caller's responsibility to close the io.ReadCloser
func (c *DropletClient) Download(guid string) (io.ReadCloser, error) {
	// This is the initial request, which will redirect to the blobstore location.
	// The client will not automatically follow this redirect and uses a secondary
	// unauthenticated client to download the bits
	// https://v3-apidocs.cloudfoundry.org/version/3.127.0/index.html#download-droplet-bits
	p := path.Format("/v3/droplets/%s/download", guid)
	req := http.NewRequest("GET", p).WithFollowRedirects(false)
	resp, err := c.client.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error getting %s: %w", p, err)
	}
	if !http.IsResponseRedirect(resp.StatusCode) {
		return nil, fmt.Errorf("error downloading droplet %s bits, expected redirect to blobstore", guid)
	}

	// get the full URL to the blobstore via the Location header
	blobStoreLocation := resp.Header.Get("Location")
	if blobStoreLocation == "" {
		return nil, errors.New("response redirect Location header was empty")
	}

	// directly download the bits from blobstore using an unauthenticated client as
	// some blob stores will return a 400 if an Authorization header is sent
	req = http.NewRequest("GET", "")
	blobstoreHTTPExecutor := http.NewExecutor(
		c.client.unauthenticatedClientProvider, blobStoreLocation, c.client.config.UserAgent)

	resp, err = blobstoreHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error downloading droplet %s bits from blobstore", guid)
	}

	return resp.Body, nil
}

// Get retrieves the droplet by ID
func (c *DropletClient) Get(guid string) (*resource.Droplet, error) {
	var d resource.Droplet
	err := c.client.get(path.Format("/v3/droplets/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// List pages all droplets the user has access to
func (c *DropletClient) List(opts *DropletListOptions) ([]*resource.Droplet, *Pager, error) {
	if opts == nil {
		opts = NewDropletListOptions()
	}
	var res resource.DropletList
	err := c.client.get(path.Format("/v3/droplets?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all droplets the user has access to
func (c *DropletClient) ListAll(opts *DropletListOptions) ([]*resource.Droplet, error) {
	if opts == nil {
		opts = NewDropletListOptions()
	}
	return AutoPage[*DropletListOptions, *resource.Droplet](opts, func(opts *DropletListOptions) ([]*resource.Droplet, *Pager, error) {
		return c.List(opts)
	})
}

// ListForApp pages all droplets for the specified app
func (c *DropletClient) ListForApp(appGUID string, opts *DropletAppListOptions) ([]*resource.Droplet, *Pager, error) {
	if opts == nil {
		opts = NewDropletAppListOptions()
	}
	var res resource.DropletList
	err := c.client.get(path.Format("/v3/apps/%s/droplets?%s", appGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppAll retrieves all droplets for the specified app
func (c *DropletClient) ListForAppAll(appGUID string, opts *DropletAppListOptions) ([]*resource.Droplet, error) {
	if opts == nil {
		opts = NewDropletAppListOptions()
	}
	return AutoPage[*DropletAppListOptions, *resource.Droplet](opts, func(opts *DropletAppListOptions) ([]*resource.Droplet, *Pager, error) {
		return c.ListForApp(appGUID, opts)
	})
}

// ListForPackage pages all droplets for the specified package
func (c *DropletClient) ListForPackage(packageGUID string, opts *DropletPackageListOptions) ([]*resource.Droplet, *Pager, error) {
	if opts == nil {
		opts = NewDropletPackageListOptions()
	}
	var res resource.DropletList
	err := c.client.get(path.Format("/v3/packages/%s/droplets?%s", packageGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForPackageAll retrieves all droplets for the specified package
func (c *DropletClient) ListForPackageAll(packageGUID string, opts *DropletPackageListOptions) ([]*resource.Droplet, error) {
	if opts == nil {
		opts = NewDropletPackageListOptions()
	}
	return AutoPage[*DropletPackageListOptions, *resource.Droplet](opts, func(opts *DropletPackageListOptions) ([]*resource.Droplet, *Pager, error) {
		return c.ListForPackage(packageGUID, opts)
	})
}

// GetCurrentAssociationForApp retrieves the current droplet relationship for an app
func (c *DropletClient) GetCurrentAssociationForApp(appGUID string) (*resource.DropletCurrent, error) {
	var d resource.DropletCurrent
	err := c.client.get(path.Format("/v3/apps/%s/relationships/current_droplet", appGUID), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetCurrentForApp retrieves the current droplet for an app
func (c *DropletClient) GetCurrentForApp(appGUID string) (*resource.Droplet, error) {
	var d resource.Droplet
	err := c.client.get(path.Format("/v3/apps/%s/droplets/current", appGUID), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// SetCurrentAssociationForApp sets the current droplet for an app. The current droplet is the droplet that the app will use when running
func (c *DropletClient) SetCurrentAssociationForApp(appGUID, dropletGUID string) (*resource.DropletCurrent, error) {
	var d resource.DropletCurrent
	r := resource.ToOneRelationship{Data: &resource.Relationship{GUID: dropletGUID}}
	_, err := c.client.patch(path.Format("/v3/apps/%s/relationships/current_droplet", appGUID), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Update an existing droplet
func (c *DropletClient) Update(guid string, r *resource.DropletUpdate) (*resource.Droplet, error) {
	var d resource.Droplet
	_, err := c.client.patch(path.Format("/v3/droplets/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
