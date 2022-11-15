package client

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"io"
	"mime/multipart"
	http2 "net/http"
	"net/url"
	"os"
)

type PackageClient commonClient

// PackageListOptions list filters
type PackageListOptions struct {
	*ListOptions

	GUIDs  Filter `qs:"guids"`  // list of package guids to filter by
	States Filter `qs:"states"` // list of package states to filter by
	Types  Filter `qs:"types"`  // list of package types to filter by, docker or bits
}

// NewPackageListOptions creates new options to pass to list
func NewPackageListOptions() *PackageListOptions {
	return &PackageListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o PackageListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Copy the bits of a source package to a target package
func (c *PackageClient) Copy(srcPackageGUID string, destAppGUID string) (*resource.Package, error) {
	var d resource.Package
	r := resource.NewPackageCopy(destAppGUID)
	_, err := c.client.post(path("/v3/packages?source_guid=%s", srcPackageGUID), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Create a new package
func (c *PackageClient) Create(r *resource.PackageCreate) (*resource.Package, error) {
	var p resource.Package
	_, err := c.client.post("/v3/packages", r, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Delete the specified package
func (c *PackageClient) Delete(guid string) error {
	_, err := c.client.delete(path("/v3/packages/%s", guid))
	return err
}

// Download the bits of an existing package
// It is the caller's responsibility to close the io.ReadCloser
func (c *PackageClient) Download(guid string) (io.ReadCloser, error) {
	// This is the initial request, which will redirect to the internal blobstore location.
	// The client should automatically follow this redirect. External blob stores are untested.
	// https://v3-apidocs.cloudfoundry.org/version/3.127.0/index.html#download-package-bits
	p := path("/v3/packages/%s/download", guid)
	req := http.NewRequest("GET", p)
	resp, err := c.client.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error getting %s: %w", p, err)
	}
	if resp.StatusCode != http2.StatusOK {
		return nil, fmt.Errorf("error getting %s, response code: %d", p, resp.StatusCode)
	}
	return resp.Body, nil
}

// Get the specified build
func (c *PackageClient) Get(guid string) (*resource.Package, error) {
	var p resource.Package
	err := c.client.get(path("/v3/packages/%s", guid), &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// List pages all the packages the user has access to
func (c *PackageClient) List(opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
	if opts == nil {
		opts = NewPackageListOptions()
	}
	var res resource.PackageList
	err := c.client.get(path("/v3/packages?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all the packages the user has access to
func (c *PackageClient) ListAll(opts *PackageListOptions) ([]*resource.Package, error) {
	if opts == nil {
		opts = NewPackageListOptions()
	}
	return AutoPage[*PackageListOptions, *resource.Package](opts, func(opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
		return c.List(opts)
	})
}

// ListForApp pages all the packages the user has access to
func (c *PackageClient) ListForApp(appGUID string, opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
	if opts == nil {
		opts = NewPackageListOptions()
	}
	var res resource.PackageList
	err := c.client.get(path("/v3/apps/%s/packages?%s", appGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppAll retrieves all the packages the user has access to
func (c *PackageClient) ListForAppAll(appGUID string, opts *PackageListOptions) ([]*resource.Package, error) {
	if opts == nil {
		opts = NewPackageListOptions()
	}
	return AutoPage[*PackageListOptions, *resource.Package](opts, func(opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
		return c.ListForApp(appGUID, opts)
	})
}

// PollReady waits until the package is ready, fails, or times out
func (c *PackageClient) PollReady(guid string, opts *PollingOptions) error {
	return PollForStateOrTimeout(func() (string, error) {
		pkg, err := c.Get(guid)
		if pkg != nil {
			return string(pkg.State), err
		}
		return "", err
	}, string(resource.PackageStateReady), opts)
}

// Update the specified attributes of the package
func (c *PackageClient) Update(guid string, r *resource.PackageUpdate) (*resource.Package, error) {
	var p resource.Package
	_, err := c.client.patch(path("/v3/packages/%s", guid), r, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UploadBits uploads an app's zip file contents
func (c *PackageClient) UploadBits(guid string, zipFile io.Reader) error {
	requestFile, err := os.CreateTemp("", "requests")
	if err != nil {
		return fmt.Errorf("could not create temp zipFile for package bits: %w", err)
	}
	defer func() {
		_ = requestFile.Close()
		_ = os.Remove(requestFile.Name())
	}()

	formWriter := multipart.NewWriter(requestFile)
	part, err := formWriter.CreateFormFile("bits", "package.zip")
	if err != nil {
		return fmt.Errorf("error uploading package %s bits: %w", guid, err)
	}

	_, err = io.Copy(part, zipFile)
	if err != nil {
		return fmt.Errorf("error uploading package %s bits, failed to copy all bytes: %w", guid, err)
	}

	err = formWriter.Close()
	if err != nil {
		return fmt.Errorf("error uploading package %s bits, failed to close multipart formWriter: %w", guid, err)
	}

	_, err = requestFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error uploading package %s bits, failed to seek beginning of temp zipFile: %w", guid, err)
	}
	fileStats, err := requestFile.Stat()
	if err != nil {
		return fmt.Errorf("error uploading package %s bits, failed to stat temp zipFile: %w", guid, err)
	}

	req := http.NewRequest("POST", path("/v3/packages/%s/upload", guid)).
		WithContentType(fmt.Sprintf("multipart/form-data; boundary=%s", formWriter.Boundary())).
		WithContentLength(fileStats.Size()).
		WithBody(requestFile)

	resp, err := c.client.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return fmt.Errorf("error uploading package %s bits: %w", guid, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http2.StatusOK {
		return c.client.handleError(resp)
	}

	return nil
}

// TODO Stage a package
