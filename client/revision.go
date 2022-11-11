package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"net/url"
)

type RevisionClient commonClient

// RevisionListOptions list filters
type RevisionListOptions struct {
	*ListOptions

	Versions Filter `qs:"versions"`
}

// NewRevisionListOptions creates new options to pass to list
func NewRevisionListOptions() *RevisionListOptions {
	return &RevisionListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o RevisionListOptions) ToQueryString() url.Values {
	return o.ListOptions.ToQueryString(o)
}

// Get the specified revision
func (c *RevisionClient) Get(guid string) (*resource.Revision, error) {
	var res resource.Revision
	err := c.client.get(path("/v3/revisions/%s", guid), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetEnvironmentVariables retrieves the specified revision's environment variables
func (c *RevisionClient) GetEnvironmentVariables(guid string) (map[string]*string, error) {
	var res resource.EnvVarResponse
	err := c.client.get(path("/v3/revisions/%s/environment_variables", guid), &res)
	if err != nil {
		return nil, err
	}
	return res.Var, nil
}

// List pages revisions that are associated with the specified app
func (c *RevisionClient) List(appGUID string, opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
	if opts == nil {
		opts = NewRevisionListOptions()
	}
	var res resource.RevisionList
	err := c.client.get(path("/v3/apps/%s/revisions?%s", appGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all revisions that are associated with the specified app
func (c *RevisionClient) ListAll(appGUID string, opts *RevisionListOptions) ([]*resource.Revision, error) {
	if opts == nil {
		opts = NewRevisionListOptions()
	}
	return AutoPage[*RevisionListOptions, *resource.Revision](opts, func(opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
		return c.List(appGUID, opts)
	})
}

// ListDeployed pages deployed revisions that are associated with the specified app
func (c *RevisionClient) ListDeployed(appGUID string, opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
	if opts == nil {
		opts = NewRevisionListOptions()
	}
	var res resource.RevisionList
	err := c.client.get(path("/v3/apps/%s/revisions/deployed?%s", appGUID, opts.ToQueryString()), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListDeployedAll pages deployed revisions that are associated with the specified app
func (c *RevisionClient) ListDeployedAll(appGUID string, opts *RevisionListOptions) ([]*resource.Revision, error) {
	if opts == nil {
		opts = NewRevisionListOptions()
	}
	return AutoPage[*RevisionListOptions, *resource.Revision](opts, func(opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
		return c.ListDeployed(appGUID, opts)
	})
}

// Update the specified attributes of the deployment
func (c *RevisionClient) Update(guid string, r *resource.RevisionUpdate) (*resource.Revision, error) {
	var res resource.Revision
	_, err := c.client.patch(path("/v3/revisions/%s", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
