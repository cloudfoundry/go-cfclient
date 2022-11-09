package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type SidecarClient commonClient

// SidecarListOptions list filters
type SidecarListOptions struct {
	*ListOptions
}

// NewSidecarListOptions creates new options to pass to list
func NewSidecarListOptions() *SidecarListOptions {
	return &SidecarListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o SidecarListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Create a new app sidecar
func (c *SidecarClient) Create(appGUID string, r *resource.SidecarCreate) (*resource.Sidecar, error) {
	var sc resource.Sidecar
	_, err := c.client.post(path("/v3/apps/%s/sidecars", appGUID), r, &sc)
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

// Delete the specified sidecar
func (c *SidecarClient) Delete(guid string) error {
	_, err := c.client.delete(path("/v3/sidecars/%s", guid))
	return err
}

// Get the specified app
func (c *SidecarClient) Get(guid string) (*resource.Sidecar, error) {
	var sc resource.Sidecar
	err := c.client.get(path("/v3/sidecars/%s", guid), &sc)
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

// ListForApp pages all sidecars associated with the specified app
func (c *SidecarClient) ListForApp(appGUID string, opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
	if opts == nil {
		opts = NewSidecarListOptions()
	}
	var res resource.SidecarList
	err := c.client.get(path("/v3/apps/%s/sidecars?%s", appGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppAll retrieves all sidecars associated with the specified app
func (c *SidecarClient) ListForAppAll(appGUID string, opts *SidecarListOptions) ([]*resource.Sidecar, error) {
	if opts == nil {
		opts = NewSidecarListOptions()
	}
	return AutoPage[*SidecarListOptions, *resource.Sidecar](opts, func(opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
		return c.ListForApp(appGUID, opts)
	})
}

// ListForProcess pages all sidecars associated with the specified process
func (c *SidecarClient) ListForProcess(processGUID string, opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
	if opts == nil {
		opts = NewSidecarListOptions()
	}
	var res resource.SidecarList
	err := c.client.get(path("/v3/processes/%s/sidecars?%s", processGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForProcessAll retrieves all sidecars associated with the specified process
func (c *SidecarClient) ListForProcessAll(processGUID string, opts *SidecarListOptions) ([]*resource.Sidecar, error) {
	if opts == nil {
		opts = NewSidecarListOptions()
	}
	return AutoPage[*SidecarListOptions, *resource.Sidecar](opts, func(opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
		return c.ListForProcess(processGUID, opts)
	})
}

// Update the specified attributes of the app
func (c *SidecarClient) Update(guid string, r *resource.SidecarUpdate) (*resource.Sidecar, error) {
	var sc resource.Sidecar
	_, err := c.client.patch(path("/v3/sidecars/%s", guid), r, &sc)
	if err != nil {
		return nil, err
	}
	return &sc, nil
}
