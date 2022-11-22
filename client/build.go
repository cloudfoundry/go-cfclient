package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type BuildClient commonClient

// BuildListOptions list filters
type BuildListOptions struct {
	*ListOptions

	States       Filter `qs:"states"`
	AppGUIDs     Filter `qs:"app_guids"`
	PackageGUIDs Filter `qs:"package_guids"`
}

// BuildAppListOptions list filters
type BuildAppListOptions struct {
	*ListOptions

	States Filter `qs:"states"`
}

// NewBuildListOptions creates new options to pass to list
func NewBuildListOptions() *BuildListOptions {
	return &BuildListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o BuildListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// NewBuildAppListOptions creates new options to pass to list
func NewBuildAppListOptions() *BuildAppListOptions {
	return &BuildAppListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o BuildAppListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new build
func (c *BuildClient) Create(ctx context.Context, r *resource.BuildCreate) (*resource.Build, error) {
	var build resource.Build
	_, err := c.client.post(ctx, "/v3/builds", r, &build)
	if err != nil {
		return nil, err
	}
	return &build, nil
}

// Delete the specified build
func (c *BuildClient) Delete(ctx context.Context, guid string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/builds/%s", guid))
	return err
}

// Get the specified build
func (c *BuildClient) Get(ctx context.Context, guid string) (*resource.Build, error) {
	var build resource.Build
	err := c.client.get(ctx, path.Format("/v3/builds/%s", guid), &build)
	if err != nil {
		return nil, err
	}
	return &build, nil
}

// List pages all builds the user has access to
func (c *BuildClient) List(ctx context.Context, opts *BuildListOptions) ([]*resource.Build, *Pager, error) {
	if opts == nil {
		opts = NewBuildListOptions()
	}
	var res resource.BuildList
	err := c.client.get(ctx, path.Format("/v3/builds?%s", opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all builds the user has access to
func (c *BuildClient) ListAll(ctx context.Context, opts *BuildListOptions) ([]*resource.Build, error) {
	if opts == nil {
		opts = NewBuildListOptions()
	}
	return AutoPage[*BuildListOptions, *resource.Build](opts, func(opts *BuildListOptions) ([]*resource.Build, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListForApp pages all builds for the app the user has access to
func (c *BuildClient) ListForApp(ctx context.Context, appGUID string, opts *BuildAppListOptions) ([]*resource.Build, *Pager, error) {
	if opts == nil {
		opts = NewBuildAppListOptions()
	}
	var res resource.BuildList
	err := c.client.get(ctx, path.Format("/v3/apps/%s/builds?%s", appGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppAll retrieves all builds for the app the user has access to
func (c *BuildClient) ListForAppAll(ctx context.Context, appGUID string, opts *BuildAppListOptions) ([]*resource.Build, error) {
	if opts == nil {
		opts = NewBuildAppListOptions()
	}
	return AutoPage[*BuildAppListOptions, *resource.Build](opts, func(opts *BuildAppListOptions) ([]*resource.Build, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// PollStaged waits until the build is staged, fails, or times out
func (c *BuildClient) PollStaged(ctx context.Context, guid string, opts *PollingOptions) error {
	return PollForStateOrTimeout(func() (string, error) {
		build, err := c.Get(ctx, guid)
		if build != nil {
			return string(build.State), err
		}
		return "", err
	}, string(resource.BuildStateStaged), opts)
}

// Update the specified attributes of the build
func (c *BuildClient) Update(ctx context.Context, guid string, r *resource.BuildUpdate) (*resource.Build, error) {
	var build resource.Build
	_, err := c.client.patch(ctx, path.Format("/v3/builds/%s", guid), r, &build)
	if err != nil {
		return nil, err
	}
	return &build, nil
}
