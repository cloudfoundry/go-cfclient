package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

// RootClient queries the global API root /
type RootClient commonClient

// Get queries the global API root /
//
// These endpoints link to other resources, endpoints, and external services that are relevant to
// authenticated API clients.
func (c *RootClient) Get(ctx context.Context) (*resource.Root, error) {
	var root resource.Root

	// NOTE - this will end up needlessly sending an auth header which the endpoint will ignore
	err := c.client.get(ctx, "/", &root)
	if err != nil {
		return nil, err
	}
	return &root, nil
}
