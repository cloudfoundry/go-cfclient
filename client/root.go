package client

import (
	"context"

	"github.com/cloudfoundry/go-cfclient/v3/resource"
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

// GetV3 queries the V3 API root /v3
//
// This endpoint returns links to all the resources available on the v3 API.
func (c *RootClient) GetV3(ctx context.Context) (*resource.V3Root, error) {
	var v3Root resource.V3Root
	// NOTE - this will end up needlessly sending an auth header which the endpoint will ignore
	if err := c.client.get(ctx, "/v3", &v3Root); err != nil {
		return nil, err
	}
	return &v3Root, nil
}
