package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type EnvVarGroupClient commonClient

// Get retrieves the specified envvar group
func (c *EnvVarGroupClient) Get(name string) (*resource.EnvVarGroup, error) {
	var e resource.EnvVarGroup
	err := c.client.get(path.Format("/v3/environment_variable_groups/%s", name), &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// GetRunning retrieves the running envvar group
func (c *EnvVarGroupClient) GetRunning() (*resource.EnvVarGroup, error) {
	return c.Get("running")
}

// GetStaging retrieves the running envvar group
func (c *EnvVarGroupClient) GetStaging() (*resource.EnvVarGroup, error) {
	return c.Get("staging")
}

// Update the specified attributes of the envar group
func (c *EnvVarGroupClient) Update(name string, r *resource.EnvVarGroupUpdate) (*resource.EnvVarGroup, error) {
	var e resource.EnvVarGroup
	_, err := c.client.patch(path.Format("/v3/environment_variable_groups/%s", name), r, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// UpdateRunning updates the specified attributes of the running envar group
func (c *EnvVarGroupClient) UpdateRunning(r *resource.EnvVarGroupUpdate) (*resource.EnvVarGroup, error) {
	return c.Update("running", r)
}

// UpdateStaging updates the specified attributes of the staging envar group
func (c *EnvVarGroupClient) UpdateStaging(r *resource.EnvVarGroupUpdate) (*resource.EnvVarGroup, error) {
	return c.Update("staging", r)
}
