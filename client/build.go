package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type BuildClient commonClient

const BuildsPath = "/v3/builds"

type BuildListOptions struct {
	*ListOptions

	States       Filter `filter:"states,omitempty"`
	AppGUIDs     Filter `filter:"app_guids,omitempty"`
	PackageGUIDs Filter `filter:"package_guids,omitempty"`
}

type BuildAppListOptions struct {
	*ListOptions

	States Filter `filter:"states,omitempty"`
}

func NewBuildListOptions() *BuildListOptions {
	return &BuildListOptions{
		ListOptions: NewListOptions(),
	}
}

func NewBuildAppListOptions() *BuildAppListOptions {
	return &BuildAppListOptions{
		ListOptions: NewListOptions(),
	}
}

func (c *BuildClient) Create(r *resource.BuildCreate) (*resource.Build, error) {
	var build resource.Build
	err := c.client.post(r.Package.GUID, BuildsPath, r, &build)
	if err != nil {
		return nil, err
	}
	return &build, nil
}

func (c *BuildClient) Delete(guid string) error {
	return c.client.delete(joinPath(BuildsPath, guid))
}

func (c *BuildClient) Get(guid string) (*resource.Build, error) {
	var build resource.Build
	err := c.client.get(joinPath(BuildsPath, guid), &build)
	if err != nil {
		return nil, err
	}
	return &build, nil
}

func (c *BuildClient) List(opts *BuildListOptions) ([]*resource.Build, *Pager, error) {
	var res resource.BuildList
	err := c.client.get(joinPathAndQS(opts.ToQueryString(opts), BuildsPath), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

func (c *BuildClient) ListAll() ([]*resource.Build, error) {
	opts := NewBuildListOptions()
	var allBuilds []*resource.Build
	for {
		builds, pager, err := c.List(opts)
		if err != nil {
			return nil, err
		}
		allBuilds = append(allBuilds, builds...)
		if !pager.HasNextPage() {
			break
		}
		opts.ListOptions = pager.NextPage(opts.ListOptions)
	}
	return allBuilds, nil
}

func (c *BuildClient) ListForApp(appGUID string, opts *BuildAppListOptions) ([]*resource.Build, *Pager, error) {
	var res resource.BuildList
	err := c.client.get(joinPathAndQS(opts.ToQueryString(opts), AppsPath, appGUID, "builds"), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

func (c *BuildClient) Update(guid string, r *resource.BuildUpdate) (*resource.Build, error) {
	var build resource.Build
	err := c.client.patch(joinPath(BuildsPath, guid), r, &build)
	if err != nil {
		return nil, err
	}
	return &build, nil
}
